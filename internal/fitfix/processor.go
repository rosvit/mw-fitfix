package fitfix

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/encoder"
	"github.com/muktihari/fit/profile/filedef"
	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/proto"
)

const garminEpoch = -631065600 * time.Second

func Process(src io.Reader, target io.Writer, opts Options) error {
	activity, err := decode(src)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	if len(activity.Records) == 0 {
		slog.Error("Failed to proceed with fit file without records")
		return errors.New("no records")
	}
	if len(activity.Sessions) == 0 {
		slog.Error("Failed to proceed with fit file without session")
		return errors.New("no session")
	}
	session := activity.Sessions[0]
	startTime := activity.Records[0].Timestamp
	endTime := activity.Records[len(activity.Records)-1].Timestamp

	fixActivity(activity, startTime, endTime)
	fixSession(session, startTime, endTime)
	activity.Events = createEvents(startTime, endTime)

	if opts.Device {
		fixFileID(&activity.FileId)
	}

	if opts.Laps {
		elapsed := session.TotalElapsedTime
		timerTime := session.TotalTimerTime
		activity.Laps = []*mesgdef.Lap{createLap(startTime, endTime, elapsed, timerTime)}
		session.SetNumLaps(1)
		if len(activity.Sessions) > 1 {
			activity.Sessions = activity.Sessions[:1]
			activity.Activity.SetNumSessions(1)
		}
	}

	return encode(activity, target)
}

func decode(r io.Reader) (*filedef.Activity, error) {
	dec := decoder.New(r)
	decoded, err := dec.Decode()
	if err != nil {
		slog.Error("Decoding failed:", KeyError, err)
		return nil, err
	}
	return filedef.NewActivity(decoded.Messages...), nil
}

func encode(activity *filedef.Activity, w io.Writer) error {
	fit := activity.ToFIT(nil)
	enc := encoder.New(w, encoder.WithProtocolVersion(proto.V2))
	if err := enc.Encode(&fit); err != nil {
		slog.Error("Encoding failed:", KeyError, err)
		return err
	}
	return nil
}

func fixActivity(a *filedef.Activity, startTime, endTime time.Time) {
	a.Activity.SetTimestamp(endTime)
	// MW writes UNIX timestamp as local timestamp => convert to Garmin Epoch and compute offset
	offset := a.Activity.LocalTimestamp.Add(garminEpoch).Sub(startTime).Round(time.Hour)
	localTimestamp := endTime.Add(offset)
	a.Activity.SetLocalTimestamp(localTimestamp)
	a.FieldDescriptions = nil
}

func fixSession(s *mesgdef.Session, startTime, endTime time.Time) {
	s.SetDeveloperFields()
	s.SetMessageIndex(0)
	s.SetStartTime(startTime)
	s.SetTimestamp(endTime)
}

func createEvents(startTime, endTime time.Time) []*mesgdef.Event {
	start := mesgdef.NewEvent(nil)
	start.SetTimestamp(startTime)
	start.SetEvent(typedef.EventTimer)
	start.SetEventType(typedef.EventTypeStart)
	start.SetEventGroup(0)

	end := mesgdef.NewEvent(nil)
	end.SetTimestamp(endTime)
	end.SetEvent(typedef.EventTimer)
	end.SetEventType(typedef.EventTypeStopAll)
	end.SetEventGroup(0)

	return []*mesgdef.Event{start, end}
}

func fixFileID(msg *mesgdef.FileId) {
	msg.SetManufacturer(typedef.ManufacturerMywhoosh)
	msg.SetProduct(0)
	msg.SetSerialNumber(0)
}

func createLap(startTime, endTime time.Time, elapsed, timerTime uint32) *mesgdef.Lap {
	return mesgdef.NewLap(nil).
		SetMessageIndex(0).
		SetTimestamp(endTime).
		SetStartTime(startTime).
		SetEvent(typedef.EventLap).
		SetEventType(typedef.EventTypeStop).
		SetSport(typedef.SportCycling).
		SetSubSport(typedef.SubSportVirtualActivity).
		SetIntensity(typedef.IntensityActive).
		SetLapTrigger(typedef.LapTriggerManual).
		SetTotalTimerTime(elapsed).
		SetTotalElapsedTime(timerTime)
}
