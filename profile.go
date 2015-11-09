package system

import (
	"math"
)

// Profile captures various statistics about the tasks of an application
// running on a platform.
type Profile struct {
	ASAP     []float64 // As Soon As Possible, the earliest start time
	ALAP     []float64 // As Late As Possible, the latest start time
	Mobility []float64 // max(0, ALAP - ASAP)

	time []float64
}

// NewProfile collects a profile of the given system. Since the mapping of
// the tasks onto the cores is assumed to be unknown at this stage, the profile
// is based on the average execution time of the tasks across all the cores.
func NewProfile(platform *Platform, application *Application) *Profile {
	nc := len(platform.Cores)
	nt := len(application.Tasks)

	profile := &Profile{
		ASAP:     make([]float64, nt),
		ALAP:     make([]float64, nt),
		Mobility: make([]float64, nt),

		time: make([]float64, nt),
	}

	for i := 0; i < nt; i++ {
		if i == 0 {
			profile.ASAP[i] = math.Inf(-1)
			profile.ALAP[i] = math.Inf(1)
		} else {
			profile.ASAP[i] = profile.ASAP[0]
			profile.ALAP[i] = profile.ALAP[0]
		}

		for j := 0; j < nc; j++ {
			profile.time[i] += platform.Cores[j].Time[application.Tasks[i].Type]
		}
		profile.time[i] /= float64(nc)
	}

	// Compute ASAP starting from the roots.
	for _, i := range application.Roots() {
		profile.propagateASAP(application, i, 0)
	}

	leafs := application.Leafs()

	totalASAP := float64(0)
	for _, i := range leafs {
		if end := profile.ASAP[i] + profile.time[i]; end > totalASAP {
			totalASAP = end
		}
	}

	// Compute ALAP starting from the leafs.
	for _, i := range leafs {
		profile.propagateALAP(application, i, totalASAP)
	}

	return profile
}

func (self *Profile) propagateASAP(application *Application, i uint, time float64) {
	if self.ASAP[i] >= time {
		return
	}

	self.ASAP[i] = time
	time += self.time[i]

	for _, i = range application.Tasks[i].Children {
		self.propagateASAP(application, i, time)
	}
}

func (self *Profile) propagateALAP(application *Application, i uint, time float64) {
	if time > self.time[i] {
		time = time - self.time[i]
	} else {
		time = 0
	}

	if time >= self.ALAP[i] {
		return
	}

	self.ALAP[i] = time

	if time > self.ASAP[i] {
		self.Mobility[i] = time - self.ASAP[i]
	} else {
		self.Mobility[i] = 0
	}

	for _, i = range application.Tasks[i].Parents {
		self.propagateALAP(application, i, time)
	}
}
