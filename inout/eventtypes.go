package inout

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	svix "github.com/svix/svix-webhooks/go"
)

func ImportEventTypesJson(sc *svix.Svix, reader io.Reader, update bool) error {
	dec := json.NewDecoder(reader)
	var eventTypes []*svix.EventTypeIn
	err := dec.Decode(&eventTypes)
	if err != nil {
		return err
	}
	for _, et := range eventTypes {
		err = createOrUpdateEventType(sc, et, update)
		if err != nil {
			return err
		}
	}
	return nil
}

func ImportEventTypesCsv(sc *svix.Svix, reader io.Reader, update bool) error {
	csvReader := csv.NewReader(reader)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) < 2 {
			return fmt.Errorf("invalid csv record")
		}
		et := &svix.EventTypeIn{
			Name:        record[0],
			Description: record[1],
		}
		err = createOrUpdateEventType(sc, et, update)
		if err != nil {
			return err
		}
	}
	return nil
}

func createOrUpdateEventType(sc *svix.Svix, et *svix.EventTypeIn, update bool) error {
	_, err := sc.EventType.Create(et)
	if err != nil {
		if sErr, ok := err.(*svix.Error); ok {
			if sErr.Status() == 409 {
				if update {
					_, err := sc.EventType.Update(et.Name, &svix.EventTypeUpdate{Description: et.Description})
					if err != nil {
						return err
					}
				}
			}
		} else {
			return err
		}
	}
	return nil
}

func GetAllEventTypes(sc *svix.Svix) ([]svix.EventTypeOut, error) {
	var eventTypes []svix.EventTypeOut
	done := false
	var iterator *string
	for !done {
		out, err := sc.EventType.List(&svix.EventTypeListOptions{
			Iterator: iterator,
		})
		if err != nil {
			return nil, err
		}
		for _, et := range out.Data {
			eventTypes = append(eventTypes, svix.EventTypeOut(et))
		}
		if out.Iterator != nil {
			iterator = out.Iterator
		}
		done = out.Done
	}
	return eventTypes, nil
}

func WriteEventTypesAsCsv(eventTypes []svix.EventTypeOut, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()
	for _, et := range eventTypes {
		err := csvWriter.Write([]string{et.Name, et.Description})
		if err != nil {
			return err
		}
	}
	return nil
}
