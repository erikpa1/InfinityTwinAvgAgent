package tools

import (
	"math"
)

type WorkLoadData struct {
	Data  interface{} `json:"date"`  // Use interface{} to allow string or number (float64)
	State interface{} `json:"state"` // Use interface{} to allow number, boolean, or NaN (float64)
}

// Define the WorkLoads struct, which uses the Data type
type WorkLoad struct {
	Name string         `json:"name"`
	Data []WorkLoadData `json:"data"`
}

// _DandDAddResult holds the result of distance and downtime calculations.
type _DandDAddResult struct {
	DowntimeAt       int64
	Downtime         int64
	DistanceFromLast float64
	WasDowntime      bool
	WasUnknown       bool
	NumberDowntimes  int64
	AvrageDowntime   float64
}

// Reset resets the _DandDAddResult fields
func (result *_DandDAddResult) Reset() {
	result.DowntimeAt = 0
	result.Downtime = 0
	result.DistanceFromLast = 0
	result.WasDowntime = false
	result.WasUnknown = false
	result.NumberDowntimes = 0
	result.AvrageDowntime = 0
}

// AnConfig is a placeholder for configuration settings
type AnConfig struct{}

// GetDowntimeTolerance returns a constant downtime tolerance
func (ac *AnConfig) GetDowntimeTolerance() int64 {
	// return int64(120000) // 1000 milliseconds as an example
	return int64(3001) // 1000 milliseconds as an example
}

// DTMath contains mathematical functions
type DTMath struct{}

// CalculateDistance calculates the distance between two points.
func (dtm *DTMath) CalculateDistance(x1, z1, x2, z2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(z2-z1, 2))
}

// DistanceAndDowntimeCounter counts distance and downtime
type DistanceAndDowntimeCounter struct {
	TimeJump int64

	Distance        float64
	Downtime        int64
	NumberDowntimes int64
	UnknownTime     int64
	UnknownDistance float64
	TotalTime       int64

	StartPeriodStr string
	EndPeriodStr   string

	LastAt      int64
	LastKnownAt int64

	LastKnownX float64
	LastKnownZ float64

	LastX            float64
	LastZ            float64
	ComesFromUnknown bool

	DowntimeTolerance int64
	DistanceThreshold float64

	SpeedCounter       int
	SpeedValuesCounter float64

	Result _DandDAddResult

	IsFirst               bool
	LastFirstX            float64
	LastFirstZ            float64
	LastFirstAt           int64
	WorkLoad              WorkLoad
	InnerDownTimeCounter  int64
	PrevWorkLoadDataState int64
	DownTimeStartDates    []int64
}

// NewDistanceAndDowntimeCounter creates a new instance of DistanceAndDowntimeCounter
func NewDistanceAndDowntimeCounter() *DistanceAndDowntimeCounter {
	ac := AnConfig{}
	return &DistanceAndDowntimeCounter{
		StartPeriodStr:    "08:00:00",
		EndPeriodStr:      "10:00:00",
		DowntimeTolerance: ac.GetDowntimeTolerance(),
		DistanceThreshold: 0.31,
		IsFirst:           true,
	}
}

func NewDistanceAndDowntimeCounterWithParams(
	downtimeTolerance int64,
	distanceThreshold float64,
) *DistanceAndDowntimeCounter {
	return &DistanceAndDowntimeCounter{
		StartPeriodStr:    "08:00:00",
		EndPeriodStr:      "10:00:00",
		DowntimeTolerance: downtimeTolerance,
		DistanceThreshold: distanceThreshold,
		IsFirst:           true,
	}
}

// SetFirst sets the first known position and time
func (self *DistanceAndDowntimeCounter) SetFirst(x, z float64, at int64) {
	self.WorkLoad = WorkLoad{}
	if !(x == 0 && z == 0) {
		self.LastX = x
		self.LastZ = z
		self.LastKnownX = x
		self.LastKnownZ = z
		self.LastKnownAt = at
		self.LastAt = at
		self.IsFirst = false
		self.LastFirstX = x
		self.LastFirstZ = z
		self.LastFirstAt = at
		self.NumberDowntimes = 0
		self.InnerDownTimeCounter = 0
		self.PrevWorkLoadDataState = -1
		self.DownTimeStartDates = make([]int64, 0)
	}
}

// Add adds a new position and time, returning the result
func (self *DistanceAndDowntimeCounter) Add(x, z float64, at int64) _DandDAddResult {
	workLoadStartData := WorkLoadData{Data: self.LastAt, State: 1}
	workFinishStartData := WorkLoadData{Data: at, State: 1}
	self.Result.Reset()
	timeDiff := at - self.LastAt

	if !(x == 0 && z == 0) {
		if self.ComesFromUnknown {
			self.ComesFromUnknown = false
			self.Result.WasUnknown = true
			self.UnknownDistance += CalculateDistance(x, z, self.LastKnownX, self.LastKnownZ)
			workLoadStartData.State = math.NaN()
			workFinishStartData.State = math.NaN()
		} else {
			timeDiffLastKnown := at - self.LastKnownAt
			timeDiffLastFirstAt := at - self.LastFirstAt
			dist := CalculateDistance(x, z, self.LastKnownX, self.LastKnownZ)
			self.Distance += dist
			workLoadStartData.State = 1
			workFinishStartData.State = 1
			distanceFromLastFirst := CalculateDistance(x, z, self.LastFirstX, self.LastFirstZ)
			if distanceFromLastFirst < self.DistanceThreshold {
				if timeDiffLastFirstAt <= self.DowntimeTolerance {
					self.InnerDownTimeCounter += timeDiffLastKnown
				}
				if timeDiffLastFirstAt > self.DowntimeTolerance {
					self.Result.WasDowntime = true
					workLoadStartData.State = 0
					workFinishStartData.State = 0
					if self.InnerDownTimeCounter != 0 {
						var localTimeDiff = self.InnerDownTimeCounter + timeDiffLastKnown
						self.Downtime += localTimeDiff
						self.Result.DowntimeAt = at - localTimeDiff
						self.Result.Downtime = localTimeDiff
						self.InnerDownTimeCounter = 0
						self.DownTimeStartDates = append(self.DownTimeStartDates, self.Result.DowntimeAt)
					} else {
						self.Downtime += timeDiffLastKnown
						self.Result.DowntimeAt = at - timeDiffLastKnown
						self.Result.Downtime = timeDiffLastKnown
					}
				} else {
					durationSec := float64(timeDiffLastKnown) / 1000
					durationHour := durationSec / 3600
					if durationHour != 0 {
						self.SpeedValuesCounter += (dist / 1000) / durationHour
						self.SpeedCounter++
					}
				}
			} else {
				durationSec := float64(timeDiffLastKnown) / 1000
				durationHour := durationSec / 3600
				if durationHour != 0 {
					self.SpeedValuesCounter += (dist / 1000) / durationHour
					self.SpeedCounter++
				}
				self.LastFirstX = x
				self.LastFirstZ = z
				self.LastFirstAt = at
			}
			self.Result.DistanceFromLast = dist
		}
		if (self.PrevWorkLoadDataState != -1 &&
			self.PrevWorkLoadDataState == int64(1) &&
			self.PrevWorkLoadDataState != int64(workLoadStartData.State.(int))) ||
			(self.PrevWorkLoadDataState == -1 && int64(workLoadStartData.State.(int)) == 0) {
			self.NumberDowntimes++
		}
		self.LastKnownAt = at
		self.LastKnownX = x
		self.LastKnownZ = z
	} else {
		self.ComesFromUnknown = true
		self.UnknownTime += timeDiff
		self.Result.WasUnknown = true
		workLoadStartData.State = math.NaN()
		workFinishStartData.State = math.NaN()
	}

	self.TotalTime += at - self.LastAt
	self.LastAt = at
	self.LastX = x
	self.LastZ = z
	self.Result.NumberDowntimes = self.NumberDowntimes
	self.Result.AvrageDowntime = float64(self.Downtime) / float64(self.NumberDowntimes)
	if self.WorkLoad.Data == nil {
		self.WorkLoad.Data = make([]WorkLoadData, 0)
	}
	lenWorkloads := len(self.WorkLoad.Data)
	if lenWorkloads > 2 && (int64(self.WorkLoad.Data[lenWorkloads-1].State.(int)) == int64(workFinishStartData.State.(int)) && int64(self.WorkLoad.Data[lenWorkloads-2].State.(int)) == int64(workFinishStartData.State.(int))) {
		self.WorkLoad.Data[lenWorkloads-1].Data = workFinishStartData.Data
	} else {
		self.WorkLoad.Data = append(self.WorkLoad.Data, workLoadStartData)
		self.WorkLoad.Data = append(self.WorkLoad.Data, workFinishStartData)
	}
	val, _ := workFinishStartData.State.(int)
	self.PrevWorkLoadDataState = int64(val)
	return self.Result
}

func (self *DistanceAndDowntimeCounter) SetCorectStartWorkloadDowntime() WorkLoad {
	// // Shift dates for state transitions to corect Downtimes
	var workLoad = WorkLoad{
		Name: self.WorkLoad.Name,
		Data: self.WorkLoad.Data,
	}
	if self.NumberDowntimes > 0 {
		for i := 1; i < len(workLoad.Data); i++ {
			if workLoad.Data[i].State != workLoad.Data[i-1].State && workLoad.Data[i-1].State == 1 {
				// Get the correct date from DownTimeStartDates
				if len(self.DownTimeStartDates) > 0 {
					correctDate := self.DownTimeStartDates[0]
					self.DownTimeStartDates = self.DownTimeStartDates[1:]
					workLoad.Data[i-1].Data = correctDate
					workLoad.Data[i].Data = correctDate
				}
			}
		}
		self.WorkLoad = workLoad
	}

	return self.WorkLoad
}

func (self *DistanceAndDowntimeCounter) ReduceWorkLoadIfexist() WorkLoad {
	var workLoad = WorkLoad{
		Name: self.WorkLoad.Name,
		Data: make([]WorkLoadData, 0),
	}
	if len(self.WorkLoad.Data) == 0 {
		return workLoad
	}
	workLoad.Data = append(workLoad.Data, self.WorkLoad.Data[0])
	for i := 1; i < len(self.WorkLoad.Data); i++ {
		if self.WorkLoad.Data[i].State != self.WorkLoad.Data[i-1].State {
			workLoad.Data = append(workLoad.Data, self.WorkLoad.Data[i-1])
			workLoad.Data = append(workLoad.Data, self.WorkLoad.Data[i])
		}
	}
	if len(self.WorkLoad.Data) > 1 && workLoad.Data[len(workLoad.Data)-1] != self.WorkLoad.Data[len(self.WorkLoad.Data)-1] {
		workLoad.Data = append(workLoad.Data, self.WorkLoad.Data[len(self.WorkLoad.Data)-1])
	}
	self.SetCorectStartWorkloadDowntime()
	return self.WorkLoad
}

func (self *DistanceAndDowntimeCounter) AddStartFinishArrangeToWorkLoad(twinName string, CustomTo int64, CustomFrom int64) WorkLoad {
	self.WorkLoad.Name = twinName
	// Add Unknown time from start to end of setuped period
	self.UnknownTime += int64(math.Abs(float64(self.LastAt - CustomTo)))
	self.UnknownTime += int64(math.Abs(float64(self.LastFirstAt - CustomFrom)))
	// Add periods with unknown time before the process starts and after the process ends to all arranged periods
	if len(self.WorkLoad.Data) > 0 {
		unknownWorkLoadSSData := WorkLoadData{Data: CustomFrom, State: "NaN"}
		unknownWorkLoadSFData := WorkLoadData{Data: self.WorkLoad.Data[0].Data, State: "NaN"}
		unknownWorkLoadFSinishData := WorkLoadData{Data: self.WorkLoad.Data[len(self.WorkLoad.Data)-1].Data, State: "NaN"}
		unknownWorkLoadFFinishData := WorkLoadData{Data: CustomTo, State: "NaN"}
		self.WorkLoad.Data = append([]WorkLoadData{unknownWorkLoadSSData, unknownWorkLoadSFData}, self.WorkLoad.Data...)
		self.WorkLoad.Data = append(self.WorkLoad.Data, unknownWorkLoadFSinishData)
		self.WorkLoad.Data = append(self.WorkLoad.Data, unknownWorkLoadFFinishData)
	}
	return self.WorkLoad
}

// ResultToJson converts the results to JSON format
func (self *DistanceAndDowntimeCounter) ResultToJson() map[string]interface{} {
	avgSpeed := 0.0
	if self.SpeedCounter > 0 {
		avgSpeed = self.SpeedValuesCounter / float64(self.SpeedCounter)
	}

	result := map[string]interface{}{
		"distance":         self.Distance,
		"downtime":         self.Downtime,
		"unknown_time":     self.UnknownTime,
		"unknown_distance": self.UnknownDistance,
		"total_time":       self.TotalTime,
		"avg_speed":        avgSpeed,
	}

	return result
}
