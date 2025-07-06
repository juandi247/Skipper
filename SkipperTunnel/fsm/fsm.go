package fsm



type fsmFunc func() fsmFunc



// initial

func ValidateFlags() fsmFunc{


	return nil
}