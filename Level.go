package fLogSystem

// Logger分級
const (
	LOGLEVELTrace    LOGLEVEL = iota // 包含最詳細訊息的記錄。 這些訊息可能包含敏感性應用程式資料。 這些訊息預設會停用，且永遠不應在生產環境中啟用
	LOGLEVELDebug                    // 開發期間用於互動式調查的記錄。 這些記錄主要應包含適用於偵錯的資訊，且不具備任何長期值
	LOGLEVELInfo                     // 追蹤應用程式一般流程的記錄。 這些記錄應具備長期值
	LOGLEVELSuccess                  // 執行時期功能或操作的完成，用來指示本程式的進程或階段任務
	LOGLEVELWarning                  // 醒目提示應用程式流程中異常或未預期事件的記錄，這些異常或未預期事件不會造成應用程式執行停止
	LOGLEVELError                    // 在目前執行流程因失敗而停止時進行醒目提示的記錄。 這些記錄應指出目前活動中的失敗，而非整個應用程式的失敗
	LOGLEVELCritical                 // 描述無法復原的應用程式或系統損毀，或需要立即注意重大失敗的記錄
	LOGLEVELNum
)

type LOGLEVEL uint32

var LOGTypeName = [7]string{"Trace", "Debug", "Info", "Success", "Warning", "Error", "Critical"}
var LOGTypeNameShot = [7]string{"TRC", "DBG", "IFO", "SCS", "WRN", "ERR", "CRT"}
