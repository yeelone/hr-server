package errno

var (
	// Common errors
	OK                  = &Errno{Code: 200, Message: "OK"}
	StatusUnauthorized  = &Errno{Code: 403, Message: "StatusUnauthorized"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	ErrValidation  = &Errno{Code: 10003, Message: "Validation failed."}
	ErrDatabase    = &Errno{Code: 10004, Message: "Database error."}
	ErrToken       = &Errno{Code: 10005, Message: "Error occurred while signing the JSON web token."}
	ErrQueryParams = &Errno{Code: 10006, Message: "Error http request params ."}

	// user errors
	ErrEncrypt           = &Errno{Code: 10101, Message: "Error occurred while encrypting the user password."}
	ErrUserNotFound      = &Errno{Code: 10102, Message: "The user was not found."}
	ErrTokenInvalid      = &Errno{Code: 10103, Message: "The token was invalid."}
	ErrPasswordIncorrect = &Errno{Code: 10104, Message: "The password was incorrect."}

	// profile errors
	ErrProfileNotFound = &Errno{Code: 10201, Message: "The profile cannot found."}
	ErrFreezeProfile   = &Errno{Code: 10202, Message: "The profile cannot freeze."}
	ErrUnFreezeProfile = &Errno{Code: 10203, Message: "The profile cannot unfreeze."}

	// group errors
	ErrGroupNotFound = &Errno{Code: 10301, Message: "The group was not found."}
	ErrLockGroup     = &Errno{Code: 10302, Message: "The group cannot lock."}
	ErrUnLockGroup   = &Errno{Code: 10303, Message: "The group cannot unlock."}
	ErrInvalidGroup  = &Errno{Code: 10304, Message: "The group cannot invalid."}
	ErrValidGroup    = &Errno{Code: 10305, Message: "The group cannot valid."}

	//template errors
	ErrTemplateInvalid       = &Errno{Code: 10401, Message: "The template was invalid,please check the file"}
	ErrUploadFileTypeInvalid = &Errno{Code: 10402, Message: "The file type was invalid,please upload the file again"}
	ErrListTemplate          = &Errno{Code: 10403, Message: "Cannot query template list"}
	ErrGetTemplate           = &Errno{Code: 10404, Message: "Cannot query template"}
	ErrGetTemplateAccount    = &Errno{Code: 10405, Message: "Cannot query template account"}
	ErrCreateTemplateAccount = &Errno{Code: 10406, Message: "Cannot create template account"}
	ErrCreateTemplate = &Errno{Code: 10407, Message: "Cannot create template"}

	//tag errors
	ErrTagNoFount = &Errno{Code: 10501, Message: "Tag not found"}

	//upload errors
	ErrUploadFormatInvalid = &Errno{Code: 10601, Message: "Upload file format error ,only support csv file "}
	ErrImport              = &Errno{Code: 10602, Message: "Import file to database err  "}
	ErrCreateFile          = &Errno{Code: 10603, Message: "Create file failed "}
	ErrSalaryCalculate     = &Errno{Code: 10604, Message: "Calculate error occur"}
	ErrSalaryProfileDetail = &Errno{Code: 10605, Message: "获取指定用户指定时间的工资明细出现错误"}
	//role and permission
	ErrUserBelongRoles = &Errno{Code: 10702, Message: "user does not belong any roles"}

	//statistics
	ErrAnnulIncome = &Errno{Code: 10801, Message: "fetch annual income error "}

	ErrWriteExcel = &Errno{Code: 10901, Message: "write data to excel failed "}
)
