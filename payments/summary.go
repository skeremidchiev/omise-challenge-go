package payments

import (
	"fmt"
	"sync"
)

type Summary struct {
	*sync.Mutex

	// failed = total - successful
	// donations
	donationsCount  int
	donationsOk     int

	// amount
	total       int
	successful  int
	average     float64

	top         [3]PaymentResult
}

func GetSummaryObject() *Summary {
	return &Summary{ &sync.Mutex{}, 0, 0, 0, 0, 0.0, [3]PaymentResult{} }
}

func (s *Summary) Add(p *PaymentResult) {
	s.Lock()

	s.donationsCount++
	s.total += p.Amount
	if p.Status == true {
		s.successful += p.Amount
		s.donationsOk++

		fnumber := float64(s.donationsCount)
		fAmount := float64(p.Amount)
		s.average = s.average * (fnumber-1) / fnumber + fAmount / fnumber

		// keeping only top 3
		for idx, data := range s.top {
			if data.Amount < p.Amount {
				switch idx {
				case 0:
					s.top[2] = s.top[1]
					s.top[1] = s.top[0]
					s.top[0] = *p
				case 1:
					s.top[2] = s.top[1]
					s.top[1] = *p
				case 2:
					s.top[2] = *p
				}

				break
			}
		}
	}

	s.Unlock()
}

func (s *Summary) GetSummary() string {
	str := fmt.Sprintf("\tDonations:\n")
	str += fmt.Sprintf("\t\tTotal:  %d\n", s.donationsCount)
	str += fmt.Sprintf("\t\tOK:     %d\n", s.donationsOk)
	str += fmt.Sprintf("\t\tFailed: %d\n", s.donationsCount - s.donationsOk)
	str += fmt.Sprintf("\tTotal:      %d\n", s.total)
	str += fmt.Sprintf("\tSuccessful: %d\n", s.successful)
	str += fmt.Sprintf("\tFailed:     %d\n", s.total - s.successful)
	str += fmt.Sprintf("\tAverage:    %f\n", s.average)

	str += fmt.Sprintf("\tTop:\n")
	for idx, d := range s.top {
		str += fmt.Sprintf("\t%d\t%s\n", idx + 1, d.Name)
	}

	return str
}