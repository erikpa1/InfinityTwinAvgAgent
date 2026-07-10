package tools

const CHART_FORMAT_TIME_EXPR = "time_expr"
const CHART_FORMAT_KMH = "kmh"
const CHART_FORMAT_KM = "km"

// ReportSection represents a section of a report
type ReportSection struct {
	Name            string                 `json:"name"`
	Headings        []string               `json:"headings"`
	Data            [][]string             `json:"data"`
	ExcelNewPage    bool                   `json:"excel_new_page"`
	Chart           string                 `json:"chart"`
	ChartUnits      string                 `json:"chart_units"`
	ChartFormatting string                 `json:"chart_formatting"`
	Translate       bool                   `json:"translate"`
	Metadata        map[string]interface{} `json:"metadata"`
	ChartData       [][]interface{}        `json:"chart_data"` // []{[string, float64]}
}

func (self *ReportSection) AddData(data ...string) {
	self.Data = append(self.Data, data)
}

func (self *ReportSection) AddChartData(data ...any) {
	self.ChartData = append(self.ChartData, data)
}

// NewReportSection creates a new ReportSection
func NewReportSection() *ReportSection {
	return &ReportSection{
		Headings:        []string{},
		Data:            [][]string{},
		ExcelNewPage:    false,
		Chart:           "",
		ChartUnits:      "",
		ChartFormatting: "",
		Translate:       true,
		Metadata:        make(map[string]interface{}),
		ChartData:       [][]interface{}{},
	}
}

// ReportStructure holds multiple sections and manages the report structure.
type ReportStructure struct {
	ReportName    string           `json:"report_name"`
	ActiveSection *ReportSection   `json:"active_section"`
	Sections      []*ReportSection `json:"sections"`
}

// NewReportStructure creates a new instance of ReportStructure.
func NewReportStructure() *ReportStructure {
	return &ReportStructure{
		ReportName:    "",
		ActiveSection: nil,
		Sections:      []*ReportSection{},
	}
}

// StartSection initializes a new section with headings.
func (rs *ReportStructure) StartSection(headings ...string) *ReportSection {
	section := NewReportSection()
	section.Headings = headings
	rs.Sections = append(rs.Sections, section)
	rs.ActiveSection = section
	return section
}

// AddSectionData adds data to the current active section.
func (rs *ReportStructure) AddSectionData(data ...string) {
	if rs.ActiveSection == nil {
		rs.StartSection() // Start a new section if no active section exists.
	}
	rs.ActiveSection.AddData(data...)

}
