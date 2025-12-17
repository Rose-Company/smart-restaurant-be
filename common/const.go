package common

import "time"

const (
	SUCCESS_STATUS           = 200
	BAD_REQUEST_STATUS       = 400
	SERVER_ERROR_STATUS      = 500
	UNAUTHORIZED_STATUS      = 401
	PERMISSION_DENIED_STATUS = 403
	NOT_FOUND_STATUS         = 404
	INTERNAL_SERVER_ERR      = 500
)

const (
	PREFIX_MAIN_POSTGRES       = "MAIN_POSTGRES"
	PREFIX_YOUPASS_DO_STORAGE  = "YOUPASS_DO_STORAGE"
	PREFIX_CRONJOB_AUTO_CANCEL = "CRONJOB_AUTO_CANCEL"
)

const ( //must NOT edit this
	ENV_GIN_DEBUG  = "GIN_DEBUG"
	ENV_RABBIT_URI = "RABBIT"
)

const (
	ENVJWTSecretKey = "JWT__SECRET_KEY"
)

var (
	DATETIME_WITH_TIMEZONE = time.RFC3339
)

const (
	USER_JWT_KEY = "USER_JWT_PROFILE"
	UserId       = "user_id"
)

const (
	POSTGRES_TABLE_NAME_COMMENT                         = "public.comment_v2"
	POSTGRES_TABLE_NAME_USER_PROFILE                    = "public.directus_users"
	POSTGRES_TABLE_NAME_USER_CLASS                      = "public.class_directus_users_2"
	POSTGRES_TABLE_NAME_CLASS                           = "public.class"
	POSTGRES_TABLE_NAME_DIRECTUS_ROLES                  = "public.directus_roles"
	POSTGRES_TABLE_NAME_VOTE_TRANSACTION                = "public.vote_transaction"
	POSTGRES_TABLE_NAME_COURSE                          = "public.course"
	POSTGRES_TABLE_NAME_QUIZ                            = "public.quiz"
	POSTGRES_TABLE_NAME_LESSON                          = "public.lesson"
	POSTGRES_TABLE_NAME_SECTION                         = "public.section"
	POSTGRES_TABLE_NAME_SECTION_PART                    = "public.section_parts"
	POSTGRES_TABLE_NAME_EXTEND_USER_PROFILE             = "public.user_profile"
	POSTGRES_TABLE_NAME_TEACHING_PROGRESS               = "public.class_directus_users_3"
	POSTGRES_TABLE_NAME_CUSTOMER_FEEDBACK               = "public.customer_feedback"
	POSTGRES_TABLE_NAME_USER_EXTENED_PROFILE            = "public.user_profile"
	POSTGRES_TABLE_NAME_REGISTRATRION_REQUEST           = "public.registration_request"
	POSTGRES_TABLE_NAME_QUESTION                        = "public.question"
	POSTGRES_TABLE_NAME_SUCCESS_QUIZ_LOG                = "public.success_quiz_log"
	POSTGRES_TABLE_NAME_SUCCESS_PART                    = "public.part"
	POSTGRES_TABLE_NAME_ANSWER                          = "public.answer"
	POSTGRES_TABLE_NAME_STUDENT_TARGET                  = "public.student_target"
	POSTGRES_TABLE_NAME_SECTION_TOPIC                   = "public.section_topic"
	POSTGRES_TABLE_NAME_AI_ANSWER                       = "public.ai_answer"
	POSTGRES_TABLE_NAME_TAG_SEARCH                      = "public.tag_search"
	POSTGRES_TABLE_NAME_QUIZ_TAG_SEARCH                 = "public.quiz_tag_search"
	POSTGRES_TABLE_NAME_PART                            = "public.part"
	POSTGRES_TABLE_NAME_EXPLANATION                     = "public.explanation"
	POSTGRES_TABLE_NAME_QUIZ_SKILL                      = "public.type"
	POSTGRES_TABLE_NAME_NOTIFICATION                    = "public.notification"
	POSTGRES_TABLE_NAME_BANNER                          = "public.banner"
	POSTGRES_TABLE_NAME_MASTER_CONFIG                   = "public.master_config"
	POSTGRES_TABLE_NAME_VOCAB                           = "public.vocab"
	POSTGRES_TABLE_NAME_VOCAB_LINKING                   = "public.vocab_linking"
	POSTGRES_TABLE_NAME_MOCK_TEST                       = "public.mock_test"
	POSTGRES_TABLE_NAME_QUIZ_PART                       = "public.quiz_part"
	POSTGRES_TABLE_STUDENT_MOCK_TEST_PROGRESS           = "public.student_mock_test_progress"
	POSTGRES_TABLE_NAME_INSTRUCTION                     = "public.instructions"
	POSTGRES_TABLE_NAME_ANSWER_COMMENT                  = "public.answer_comment"
	POSTGRES_TABLE_NAME_REVIEW                          = "public.review"
	POSTGRES_TABLE_NAME_SHIFT                           = "public.shifts"
	POSTGRES_TABLE_NAME_ROOM                            = "public.room"
	POSTGRES_TABLE_CLASS_LESSON                         = "public.class_lesson"
	POSTGRES_TABLE_NAME_ANSWER_REVIEW_CONFIG            = "public.answer_review_config"
	POSTGRES_TABLE_INSTRUCTION                          = "public.instructions"
	POSTGRES_TABLE_NAME_FEEDBACK                        = "public.feedback"
	POSTGRES_TABLE_NAME_ACTIVITY_LOG                    = "public.activity_log"
	POSTGRES_TABLE_NAME_REGISTRATION_REQUEST_TAG_SEARCH = "public.registration_request_tag_search"
	POSTGRES_TABLE_NAME_LOGICAL_FRAMEWORK               = "public.writing_logic_frame"
	POSTGRES_TABLE_NAME_AI_USAGE_COUNT                  = "public.ai_usage_count"
	POSTGRES_TABLE_NAME_IELTS_BLOG                      = "public.ielts_blog"
	POSTGRES_TABLE_NAME_TAG_POSITION                    = "public.tag_position"
	POSTGRES_TABLE_NAME_TAG_POSITION_TAG_SEARCH         = "public.tag_position_tag_search"
	POSTGRES_TABLE_NAME_AI_USAGE_COUNT_TRANSACTION      = "public.ai_usage_transaction"
	POSTGRES_TABLE_NAME_AI_USAGE_TRANSACTION            = "public.ai_usage_transaction"
	POSTGRES_TABLE_NAME_EMAIL_TRANSACTION               = "public.email_transaction"
	POSTGRES_TABLE_NAME_USER_VOCAB_CATEGORY             = "public.user_vocab_category"
	POSTGRES_TABLE_NAME_VOUCHER_CODE_USER_STATS         = "public.voucher_code_user_stats"
	POSTGRES_TABLE_NAME_VOUCHER_CODES                   = "public.voucher_codes"
	POSTGRES_TABLE_NAME_VOUCHER_CODES_USER              = "voucher_code_users"
	POSTGRES_TABLE_NAME_REVIEW_PROGRESS                 = "review_progress"
	POSTGRES_TABLE_NAME_USER_VOCAB_BANK                 = "public.user_vocab_bank"
	POSTGRES_TABLE_NAME_DIRECTUS_SESSION                = "public.directus_sessions"
	POSTGRES_TABLE_NAME_TRACKING                        = "public.tracking"
)

const (
	STUDENT_LEARNED_TIME_MIN_INTERVAL int = 1 // minutes
	STUDENT_LEARNED_TIME_INTERVAL     int = 3 // minutes
)

const (
	StatusPublished             = "published"
	CollectionQuiz              = "quiz"
	AnswerStatisticByQuiz       = 1
	AnswerStatisticQuestionType = 2
	AnswerStatisticByPassage    = 3

	QuizTypeUnknown     = 0
	QuizTypeExercise    = 1
	QuizTypeAssignment  = 2
	QuizTypeTest        = 3
	QuizTypeMockTest    = 4
	VocabLevelWord      = 3
	VocabLevelSentence  = 2
	VocalLevelParagraph = 1

	MockTestTypeUnknown = 0
	MockTestTypeFull    = 1
	MockTestTypePart    = 2 // by section or passage

	BandScore0   = 0
	BandScore1   = 1
	BandScore2   = 2
	BandScore2_5 = 2.5
	BandScore3   = 3
	BandScore3_5 = 3.5
	BandScore4   = 4
	BandScore4_5 = 4.5
	BandScore5   = 5
	BandScore5_5 = 5.5
	BandScore6   = 6
	BandScore6_5 = 6.5
	BandScore7   = 7
	BandScore7_5 = 7.5
	BandScore8   = 8
	BandScore8_5 = 8.5
	BandScore9   = 9

	MinCorrectQuesBand0   = 0
	MinCorrectQuesBand1   = 1
	MinCorrectQuesBand2   = 2
	MinCorrectQuesBand2_5 = 3
	MinCorrectQuesBand3   = 5
	MinCorrectQuesBand3_5 = 7
	MinCorrectQuesBand4   = 11
	MinCorrectQuesBand4_5 = 13
	MinCorrectQuesBand5   = 16
	MinCorrectQuesBand5_5 = 18
	MinCorrectQuesBand6   = 23
	MinCorrectQuesBand6_5 = 26
	MinCorrectQuesBand7   = 30
	MinCorrectQuesBand7_5 = 32
	MinCorrectQuesBand8   = 35
	MinCorrectQuesBand8_5 = 37
	MinCorrectQuesBand9   = 39

	Part1                    = "part_1"
	Part2                    = "part_2"
	Part3                    = "part_3"
	Part4                    = "part_4"
	TotalPartsReadingSkill   = 3
	TotalPartsListeningSkill = 4

	QuizSkillTypeWriting              = 3
	QuizSkillTypeSpeaking             = 4
	QuizSkillTypeWritingSelfPractice1 = 5
	QuizSkillTypeWritingSelfPractice2 = 7

	ConfigMockTestBookCode = "MOCK_TEST_BOOK"
	QuizFullAlias          = "full"

	QuizSubmittedStatusUnknown = 0
	QuizSubmittedStatusYes     = 1
	QuizSubmittedStatusNo      = 2

	AnswerStatusDraft     = "draft"
	AnswerStatusReviewed  = "reviewed"
	AnswerStatusCompleted = "completed"

	StatusInactive = 0
)

const (
	ClaudeMaxToken                       = 1024
	AnswerReviewTypeGoogleDocs           = 1
	AnswerReviewTypeAI                   = 2
	PromptingTypeLRGR                    = "LR_GR"
	PromptingTypeFlow                    = "FLOW"
	PromptingTypeTeacherGuideline        = "SUGGESTION"
	PromptingTypeUpgrade                 = "UPGRADE"
	PromptingTypeTRCC                    = "TRCC"
	PromptingTypeCheckValidWritingAnswer = "CHECK_VALID_WRITING_ANSWER"
	PromptingTypeBandScore               = "BAND_SCORE"
	PromptingTypeLRGRTask1               = "LR_GR_TASK_1"
	PromptingTypeTRCCTask1               = "TRCC_TASK_1"
	PromptingTypeBandScoreTask1          = "BAND_SCORE_TASK_1"
	PromptingTypeUgradeTask1             = "UPGRADE_TASK_1"
	AnswerErrorVocabulary                = "vocabulary"
	AnswerErrorGrammar                   = "grammar"
	AnswerErrorCoherenceCohesion         = "coherence_cohesion"
	AnswerErrorTaskResponse              = "task_response"
	AnswerErrorDefaultLRGR               = "default_lrgr"
	AnswerErrorDefaultTRCC               = "default_trcc"
	ReviewPromptingStatusInactive        = 1
	ReviewPromptingStatusInProgress      = 2
	ReviewPromptingStatusFailed          = 3
	ReviewPromptingStatusCancelled       = 4
	PromptingTypeLRGR84                  = "LR_GR_84"
)

const (
	StatusActive  = 1
	StatusLike    = 1
	StatusDislike = 2
	StatusUnVote  = 0
)

var ObjectTypeToTableName = map[string]string{
	"course":  POSTGRES_TABLE_NAME_COURSE,
	"quiz":    POSTGRES_TABLE_NAME_QUIZ,
	"lesson":  POSTGRES_TABLE_NAME_LESSON,
	"vocab":   POSTGRES_TABLE_NAME_VOCAB,
	"comment": POSTGRES_TABLE_NAME_COMMENT,
	"review":  POSTGRES_TABLE_NAME_REVIEW,
}

const (
	VocabID    = "vocab_id"
	CourseID   = "course_id"
	QuizID     = "quiz_id"
	LessonID   = "lesson_id"
	CommentID  = "comment_id"
	FeedbackID = "feedback_id"
	AnswerID   = "answer_id"
	ReviewID   = "review_id"
)

const (
	IsRetry = "is_retry"
)
const (
	EnrollClassModeEntranceTest           = "entrance-test"
	EnrollClassModePaidClass              = "paid-class"
	EnrollClassAction                     = "ENROLL_TO_CLASS"
	EnrollClassDescription                = "Học sinh tham gia vào lớp"
	SuccessEnrollMessage                  = "Successfully enrolled"
	ErrNotFoundStudentMessage             = "Not found student"
	ErrNotFoundClassMessage               = "Not found class"
	ErrNotFoundCourseMessage              = "Not found course"
	ErrNotFoundRegistrationRequestMessage = "Not found registration request"
	ErrAlreadyEnrollMessage               = "Already enrolled"
	ErrNotEntranceTestClass               = "Class doesn't allow entrance test"
	ErrUserNotMatchRegistrationRequest    = "User not match registration request profile"
	EnrollClassStatusUnknown              = 0
	EnrollClassStatusSuccess              = 1
	EnrollClassStatusFailed               = 2
	EnrollClassStatusCanceled             = 3
)

const (
	PasswordDefault = "123456"
)
const (
	ObjectTypeVocab   = "vocab"
	ObjectTypeLesson  = "lesson"
	ObjectTypeQuiz    = "quiz"
	ObjectTypeCourse  = "course"
	ObjectTypeComment = "comment"
	ObjectTypeReview  = "review"
)

const (
	MaxFileAttached = 3
)

const (
	StatusSuccessful = 0
	StatusPartial    = 2
	StatusFailed     = 1
)

const (
	ChangeAssignee = "CHANGE_ASSIGNEE"
	ChangeReviewer = "CHANGE_REVIEWER"
	AddTag         = "ADD_TAG"
)

const (
	NotificationCommentTitle  = "Có bình luận mới kìa"
	DefaultAttachmentsComment = "tệp tin"
)

const (
	LarkMessageTypeText        = "text"
	LarkMessageTypePost        = "post"
	LarkMessageTypeInteractive = "interactive"
	LarkTagDiv                 = "div"
	LarkTagText                = "text"
	LarkTagMd                  = "lark_md"
	LarkTagPlainText           = "plain_text"
	LarkTagButton              = "button"
	LarkActionTypeDefault      = "default"
	LarkLanguageEnUs           = "en-us"
	LarkTagLink                = "a"
)

const (
	UserRoleEndUser = "END_USER"
)

const (
	CategoryWritingAISelfPracticeV2 = "writing_self_practice_v2"
)

const (
	PeriodTypeUnknown = 0
	PeriodTypeDay     = 1
	PeriodTypeWeek    = 2
	PeriodTypeMonth   = 3
	PeriodTypeYear    = 4
)

const (
	SystemName = "SYSTEM"
)

const (
	SuggestPhrase        = "SUGGEST_PHRASE"
	SourceAIFromVieToEng = "VIE_ENG"
)

const (
	AdvertisePositionWritingAIReview = "WRITING_AI_REVIEW"
)

const (
	WritingQuestionTypePosition     = "writing_question_types_search"
	WritingTask1TopicSearchPosition = "writing_task_1_topics_search"
	WritingTaskType1                = 1
	WritingTaskType2                = 2
)

const (
	AIUsageCategoryWritingSelfPractice = "writing_self_practice_v2"
	AIUsageCategoryVocabTranslate      = "vocab_translate"
)
const (
	UserVocabCategory = "USER_VOCAB_CATEGORY"
)

const (
	EmailTransactionObjectAnswer = "answer"
	ActionSubmitCorrection       = "submit_correction"
)

const (
	CommonStatusActive = 1
)

const (
	ClassStatusPublished = "published"
	ClassStatusDraft     = "draft"
)

const (
	PaymentTypeLink = "link"
)

const (
	ReviewSubProgressStatusInProcess = 1
	ReviewSubProgressStatusSuccess   = 2
	ReviewSubProgressStatusFailed    = 3
)

const (
	UserVocabCategoryKey = "USER_VOCAB_CATEGORY"
	Category             = "category"
)

const (
	UserVocabStatusInactive = 0
	UserVocabStatusLearning = 1
	UserVocabStatusLearned  = 2
)

const (
	SortStatusDesc    = "status.desc"
	SortStatusAsc     = "status.asc"
	SortCreatedAtDesc = "created_at.desc"
	SortCreatedAtAsc  = "created_at.asc"
)

const (
	NotApplicable = "N/A"
)

const (
	ReviewProgressStatusInProcess = 1
	ReviewProgressStatusSuccess   = 2
	ReviewProgressStatusFailed    = 3
)

const (
	OrderVoucherContextYoupass = 1
	OrderVoucherContext84      = 2
)

const (
	IssDirectus                = "directus"
	NanoIDAlphabet             = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	RefreshTokenModeCookie     = "cookie"
	RefreshTokenModeJson       = "json"
	DirectusRefreshTokenPrefix = "directus_refresh_token"
	DirectusProviderGoogle     = "google"
)

const (
	RegistrationStatusCancelled = "CANCELLED"
)
