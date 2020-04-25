package pdfpagedata

type PageData struct {
	Exam       ExamDetails         `json:"exam"`
	Author     AuthorDetails       `json:"author"`
	Page       PageDetails         `json:"page"`
	Contact    ContactDetails      `json:"page"`
	Questions  []QuestionDetails   `json:"questions"`
	Processing []ProcessingDetails `json:"processing"`
	Custom     []CustomDetails     `json:"custom""`
}

type ExamDetails struct {
	CourseCode string `json:"courseCode"`
	Diet       string `json:"diet"`
	UUID       string `json:"UUID"`
}

type AuthorDetails struct {
	ExamNumber string `json:"examNumber"`
	UUID       string `json:"UUID"`
}

type PageDetails struct {
	UUID   string `json:"UUID"`
	Number int    `json:"number"`
}

type ContactDetails struct {
	Name    string `json:"name"`
	UUID    string `json:"UUID"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// use section for (a), (b) and number for (i)
type QuestionDetails struct {
	UUID           string            `json:"UUID"`
	Name           string            `json:"name"` //what to call it in a dropbox etc
	Section        string            `json:"section"`
	Number         int               `json:"number"` //No Harry Potter Platform 9&3/4 questions
	Parts          []QuestionDetails `json:"parts"`
	MarksAvailable float64           `json:"marksAvailable"`
	MarksAwarded   float64           `json:"marksAwarded"`
	Marking        []MarkingAction   `json:"markers"`
	Moderating     []MarkingAction   `json:"moderators"`
	Checking       []MarkingAction   `json:"checkers"`
}

type MarkingAction struct {
	Contact  ContactDetails `json:"contact"`
	Mark     MarkDetails    `json:"mark"`
	Done     bool           `json:"done"`
	UnixTime int64          `json:"unixTime"`
	Custom   CustomDetails  `json:"contact"`
}

type MarkDetails struct {
	Given     float64 `json:"given"`
	Available float64 `json:"available"`
	Comment   float64 `json:"comment"`
}

type CustomDetails struct {
	Key   string `json:"name"`
	Value string `json:"value"`
}

type ProcessingDetails struct {
	UUID       string             `json:"UUID"`
	Previous   string             `json:"previous"`
	UnixTime   int64              `json:"unixTime"`
	Name       string             `json:"name"`
	Parameters []ParameterDetails `json:"parameters"`
	By         ContactDetails     `json:"by"`
}

type ParameterDetails struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Sequence int    `json:sequence"`
}

const (
	StartTag       = "<gradex-pagedata>"
	EndTag         = "</gradex-pagedata>"
	StartTagOffset = len(StartTag)
	EndTagOffset   = len(EndTag)
)
