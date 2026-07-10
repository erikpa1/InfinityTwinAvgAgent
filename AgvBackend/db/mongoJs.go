package db

//Created from https://claude.ai/chat/8ca9839b-0248-4290-b238-fe04b2518a31

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"turtle/credentials"
)

// VectorSearchResult represents a single search result
type VectorSearchResult struct {
	ID         string    `bson:"_id"`
	Content    string    `bson:"content"`
	Similarity float64   `bson:"similarity"`
	Embedding  []float32 `bson:"embedding,omitempty"`
}

// ClusterResult represents clustering results
type ClusterResult struct {
	Clusters  [][]VectorSearchResult `bson:"clusters"`
	Centroids [][]float32            `bson:"centroids"`
}

// VectorStats represents vector statistics
type VectorStats struct {
	TotalDocuments  int     `bson:"totalDocuments"`
	AvgVectorLength float64 `bson:"avgVectorLength"`
	MaxVectorLength int     `bson:"maxVectorLength"`
	MinVectorLength int     `bson:"minVectorLength"`
}

// InstallJavaScriptFunctions installs the vector search functions in MongoDB
func (self *AnyDBConnection) InstallJavaScriptFunctions() error {
	// JavaScript functions to install (you can add all the functions from the previous artifact)
	functions := map[string]string{
		"cosineSimilarity": `
			function cosineSimilarity(vectorA, vectorB) {
				if (!vectorA || !vectorB || vectorA.length !== vectorB.length) {
					return 0;
				}
				
				let dotProduct = 0;
				let normA = 0;
				let normB = 0;
				
				for (let i = 0; i < vectorA.length; i++) {
					dotProduct += vectorA[i] * vectorB[i];
					normA += vectorA[i] * vectorA[i];
					normB += vectorB[i] * vectorB[i];
				}
				
				if (normA === 0 || normB === 0) {
					return 0;
				}
				
				return dotProduct / (Math.sqrt(normA) * Math.sqrt(normB));
			}`,
		"vectorSearch": `
			function vectorSearch(queryVector, limit = 10, similarityThreshold = 0.7) {
				const collection = db.documents;
				let results = [];
				
				collection.find({}).forEach(function(doc) {
					if (doc.embedding && Array.isArray(doc.embedding)) {
						const similarity = cosineSimilarity(queryVector, doc.embedding);
						
						if (similarity >= similarityThreshold) {
							results.push({
								_id: doc._id,
								content: doc.content,
								similarity: similarity
							});
						}
					}
				});
				
				results.sort(function(a, b) {
					return b.similarity - a.similarity;
				});
				
				return results.slice(0, limit);
			}`,
		"findSimilarDocuments": `
			function findSimilarDocuments(documentId, limit = 5) {
				const collection = db.documents;
				const sourceDoc = collection.findOne({_id: documentId});
				
				if (!sourceDoc || !sourceDoc.embedding) {
					return {error: "Document not found or has no embedding"};
				}
				
				return vectorSearch(sourceDoc.embedding, limit + 1)
					.filter(function(doc) {
						return !doc._id.equals(documentId);
					})
					.slice(0, limit);
			}`,
		"getVectorStats": `
			function getVectorStats() {
				const collection = db.documents;
				
				return collection.aggregate([
					{
						$match: {
							embedding: {$exists: true}
						}
					},
					{
						$group: {
							_id: null,
							totalDocuments: {$sum: 1},
							avgVectorLength: {$avg: {$size: "$embedding"}},
							maxVectorLength: {$max: {$size: "$embedding"}},
							minVectorLength: {$min: {$size: "$embedding"}}
						}
					}
				]).toArray()[0];
			}`,
	}

	// Install each function
	for name, body := range functions {
		doc := bson.M{
			"_id":   name,
			"value": body,
		}

		_, err := self.Col("functions").ReplaceOne(
			context.Background(),
			bson.M{"_id": name},
			doc,
			options.Replace().SetUpsert(true),
		)
		if err != nil {
			return fmt.Errorf("failed to install function %s: %w", name, err)
		}
	}

	fmt.Println("JavaScript functions installed successfully")
	return nil
}

// VectorSearch performs vector search using MongoDB server-side JavaScript
func (self *AnyDBConnection) VectorSearch(ctx context.Context, container string, query []float32, limit int, threshold float64) ([]VectorSearchResult, error) {

	// Prepare JavaScript code to execute
	jsCode := fmt.Sprintf(`
		function() {
			const queryVector = %v;
			const limit = %d;
			const threshold = %f;
			return vectorSearch(queryVector, limit, threshold);
		}
	`, query, limit, threshold)

	// Execute the JavaScript function
	var results []VectorSearchResult

	err := self.Mongo.Database(credentials.GetDBName()).RunCommand(ctx, bson.M{
		"eval": jsCode,
	}).Decode(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to execute vector search: %w", err)
	}

	return results, nil
}
