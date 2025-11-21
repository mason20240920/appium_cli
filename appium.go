package appium_cli

import (
	"fmt"
	"strings"
	"time"

	sf "github.com/sa-/slicefunk"
)

// CreateSession Create the appium new session to do the test, Get the screen size of device
func (capModel *DeviceCapabilityModel) CreateSession() (deviceModel *DeviceDriverModel, serverErr *AppiumError) {
	// Make the struct to the desired capabilities
	model := &AppiumParameter{DesiredCapabilities: DesiredCapabilities{
		PlatformName:          capModel.Platform.ToString(),
		PlatformVersion:       capModel.PlatformVersion,
		DeviceName:            capModel.DeviceName,
		AppPackage:            capModel.AppPackage,
		AppActivity:           capModel.AppActivity,
		NewCommandTimeout:     capModel.NewCommandTimeout,
		AndroidInstallTimeout: capModel.AndroidInstallTimeout,
		AutomationName:        capModel.AutomationName,
		SystemPort:            capModel.SystemPort,
		Udid:                  capModel.Udid,
		NoReset:               capModel.NoReset,
	}}
	var result SessionResponse
	var errorResult SessionErrorResponse

	client := capModel.Client

	resp, err := client.R().
		SetBody(model).
		SetSuccessResult(&result).
		SetErrorResult(&errorResult).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session", capModel.Port))
	if err != nil {
		serverErr = &AppiumError{
			Message:   err.Error(),
			ErrorCode: CreateSessionError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   errorResult.Value.Message,
			ErrorCode: CreateSessionError,
		}
		return
	}
	deviceModel = &DeviceDriverModel{
		SessionId:  result.SessionId,
		Port:       capModel.Port,
		DeviceName: capModel.DeviceName,
		Client:     client}
	return
}

// CloseSession close the appium session
func (driver DeviceDriverModel) CloseSession() (serverErr *AppiumError) {
	resp, err := driver.Client.R().
		Delete(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Stop Session Error",
			ErrorCode: StopSessionError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Stop Session Error",
			ErrorCode: StopSessionError,
		}
		return
	}
	return
}

// FindElement
//
//	@Description: Find the element id
//	@receiver driver
//	@param param
//	@return elementId
//	@return serverErr
func (driver DeviceDriverModel) FindElement(param *FindElementPoint) (elementId string, serverErr *AppiumError) {
	var result ElementResponse

	resp, err := driver.Client.R().
		SetBody(&FindElementParam{
			Using: param.GetUsingType(),
			Value: param.Value,
		}).
		SetSuccessResult(&result).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/element", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Find Element Error",
			ErrorCode: NotFoundElement,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Create Session Error",
			ErrorCode: NotFoundElement,
		}
		return
	}
	elementId = result.Value.ELEMENT
	return
}

func (driver DeviceDriverModel) FindElements(param *FindElementPoint) (elementId []string, serverErr *AppiumError) {
	var result ElementsResponse

	resp, err := driver.Client.R().
		SetBody(&FindElementParam{
			Using: param.GetUsingType(),
			Value: param.Value,
		}).
		SetSuccessResult(&result).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/elements", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Find Element Error",
			ErrorCode: NotFoundElement,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Create Session Error",
			ErrorCode: NotFoundElement,
		}
		return
	}
	elementId = sf.Map(result.Value, func(item ElementValueModel) string {
		return item.ELEMENT
	})
	return
}

// ActionElement take action the device element
func (driver DeviceDriverModel) ActionElement(elementParam *ActionNormalParam, action ActionType) (pressTime int64, serverErr *AppiumError) {
	var result SessionResponse
	var body any
	var requestUrl string
	if action == SendKeys {
		body = &SendKeysParam{Text: elementParam.Text}
		requestUrl = fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/element/%s/value", driver.Port, driver.SessionId, elementParam.Element)
	} else if action == Clear {
		body = &SendKeysParam{Text: elementParam.Text}
		requestUrl = fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/element/%s/clear", driver.Port, driver.SessionId, elementParam.Element)
	} else {
		requestUrl = fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/element/%s/click", driver.Port, driver.SessionId, elementParam.Element)
		body = &ActionRequestParam{Element: elementParam.Element}
	}
	now := time.Now()
	// 将当前时间转换为从1970年开始的毫秒数
	pressTime = now.UnixMilli()
	resp, err := driver.Client.R().
		SetBody(body).
		SetSuccessResult(&result).
		Post(requestUrl)
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Action Element Error",
			ErrorCode: ActionElementError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Create Session Error",
			ErrorCode: ActionElementError,
		}
		return
	}
	return
}

// TouchActionByLoc make the location x and y to action
func (driver DeviceDriverModel) TouchActionByLoc(coordinate Coordinate) (serverErr *AppiumError) {
	// Create the parameters
	var actions []ActionRequestChain
	actions = append(actions, ActionRequestChain{
		Type:     "pointerMove",
		Duration: coordinate.GetDuration(),
		X:        coordinate.GetPosition().X,
		Y:        coordinate.GetPosition().Y,
	})
	actions = append(actions, ActionRequestChain{
		Type:   "pointerDown",
		Button: 0,
	})
	actions = append(actions, ActionRequestChain{
		Type:     "pause",
		Duration: 5,
	})
	actions = append(actions, ActionRequestChain{
		Type:   "pointerUp",
		Button: 0,
	})
	pointerParam := ActionRequestParams{
		PointerType: "touch",
	}
	actionParam := ActionsRequest{
		Actions:    actions,
		Parameters: pointerParam,
		Id:         "finger1",
		Type:       "pointer",
	}
	var actionTemp []ActionsRequest
	actionTemp = append(actionTemp, actionParam)
	requestParams := &ActionRequestArr{
		Actions: actionTemp,
	}
	var result SessionResponse

	resp, err := driver.Client.R().
		SetBody(requestParams).
		SetSuccessResult(&result).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/actions", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Touch Action Error",
			ErrorCode: TouchActionError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Touch Action Error",
			ErrorCode: TouchActionError,
		}
		return
	}
	return
}

// StartActivity
//
//	@Description: Start the activity of the another application
//	@receiver driver
//	@param param
//	@return serverErr
func (driver DeviceDriverModel) StartActivity(param *StartActivityParam) (serverErr *AppiumError) {
	var result SessionResponse

	resp, err := driver.Client.R().
		SetBody(param).
		SetSuccessResult(&result).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/appium/device/start_activity", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Start Activity Error",
			ErrorCode: StartActivityError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Start Activity Error",
			ErrorCode: StartActivityError,
		}
		return
	}
	return
}

// ImplicitWait
//
//	@Description: Set the amount of time the driver should wait when searching for elements
//	@receiver driver
//	@param seconds
//	@return serverErr
func (driver DeviceDriverModel) ImplicitWait(seconds time.Duration) (serverErr *AppiumError) {
	var result SessionResponse

	resp, err := driver.Client.R().
		SetBody(&ImplicitWaitParam{
			Seconds: int(seconds / time.Millisecond),
		}).
		SetSuccessResult(&result).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/timeouts/implicit_wait", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "ImplicitWait Time Out",
			ErrorCode: ImplicitWaitError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "ImplicitWait Time Out",
			ErrorCode: ImplicitWaitError,
		}
		return
	}
	return
}

// ElementActionMov
//
//		@Description: Clicks element at its center point. If the element's center point is obscured by another element,
//		an element click intercepted error is returned. If the element is outside the viewport, an element not interactable
//		error is returned. Not all drivers automatically scroll the element into view and may need to be scrolled to in order to interact with it.
//		@receiver driver
//		@param param the find way of the element
//	 @param seconds the seconds to implicitWait
//	 @param action take the action
//		@return serverErr
func (driver DeviceDriverModel) ElementActionMov(param *FindElementPoint, seconds time.Duration, action ActionType, sendKeys string) (elementId string, serverErr *AppiumError) {
	// 1. Confirm to find element
	if seconds != 0 {
		serverErr = driver.ImplicitWait(seconds)
	} else {
		serverErr = driver.ImplicitWait(5)
	}
	if serverErr != nil {
		return
	}
	// 2. Find the element
	elementId, serverErr = driver.FindElement(param)
	if serverErr != nil {
		return
	}
	// 3. Touch or move the element
	_, serverErr = driver.ActionElement(&ActionNormalParam{
		Element: elementId,
		Text:    sendKeys,
	}, action)
	return
}

// ElementActionWithTime
//
//	@Description: Element action with time
//	@receiver driver
//	@param param
//	@param seconds
//	@param action
//	@param sendKeys
//	@return elementId
//	@return serverErr
func (driver DeviceDriverModel) ElementActionWithTime(param *FindElementPoint, seconds time.Duration, action ActionType, sendKeys string) (actionTime int64, elementId string, serverErr *AppiumError) {
	// 1. Confirm to find element
	if seconds != 0 {
		serverErr = driver.ImplicitWait(seconds)
	} else {
		serverErr = driver.ImplicitWait(5)
	}
	if serverErr != nil {
		return
	}
	// 2. Find the element
	elementId, serverErr = driver.FindElement(param)
	if serverErr != nil {
		return
	}
	// 3. Touch or move the element
	actionTime, serverErr = driver.ActionElement(&ActionNormalParam{
		Element: elementId,
		Text:    sendKeys,
	}, action)
	return
}

func (driver DeviceDriverModel) WebViewElementAct(param *FindElementPoint, seconds time.Duration, action ActionType, sendKeys string) (elementId string, serverErr *AppiumError) {
	time.Sleep(seconds)
	if serverErr != nil {
		return
	}
	// 2. Find the element
	elementId, serverErr = driver.FindElement(param)
	if serverErr != nil {
		return
	}
	// 3. Touch or move the element
	_, serverErr = driver.ActionElement(&ActionNormalParam{
		Element: elementId,
		Text:    sendKeys,
	}, action)
	return
}

// GetAttribute
//
//	@Description: Get the element attribute
//	@receiver driver
//	@param param
//	@param elementId
//	@return serverErr
func (driver DeviceDriverModel) GetAttribute(param *AttributeModel, element *FindElementPoint) (value AttributeInterface, elementId string, serverErr *AppiumError) {
	// 1. Find the element
	elementId, serverErr = driver.FindElement(element)
	if serverErr != nil {
		return
	}

	var result AttributeResponse

	resp, err := driver.Client.R().
		SetSuccessResult(&result).
		Get(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/element/%s/attribute/%s", driver.Port, driver.SessionId, elementId, param.GetAttributeStr()))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Not Found the attribute",
			ErrorCode: NotFoundAttribute,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Not Found the attribute",
			ErrorCode: NotFoundAttribute,
		}
		return
	}
	value = &AttributeRetModel{Value: result.Value}
	return
}

func (driver DeviceDriverModel) GetAttributeByElementId(param *AttributeModel, elementId string) (value AttributeInterface, serverErr *AppiumError) {
	var result AttributeResponse

	resp, err := driver.Client.R().
		SetSuccessResult(&result).
		Get(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/element/%s/attribute/%s", driver.Port, driver.SessionId, elementId, param.GetAttributeStr()))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Not Found the attribute",
			ErrorCode: NotFoundAttribute,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Not Found the attribute",
			ErrorCode: NotFoundAttribute,
		}
		return
	}
	value = &AttributeRetModel{Value: result.Value}
	return
}

// TerminateApp terminate the application
func (driver DeviceDriverModel) TerminateApp(appId string) (ret bool, serverErr *AppiumError) {
	var result TerminateResponse
	resp, err := driver.Client.R().
		SetBody(&AppPropParam{
			AppId: appId,
		}).
		SetSuccessResult(&result).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/appium/device/terminate_app", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Terminate the app error",
			ErrorCode: TerminalAppError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Terminate the app error",
			ErrorCode: TerminalAppError,
		}
		return
	}
	ret = result.Value
	return
}

func (driver DeviceDriverModel) ClearText(element *FindElementPoint) (serverErr *AppiumError) {
	serverErr = driver.ImplicitWait(500 * time.Millisecond)
	if serverErr != nil {
		return
	}
	// 1. Find the element
	elementId, serverErr := driver.FindElement(element)
	if serverErr != nil {
		return
	}
	time.Sleep(300 * time.Millisecond)
	_, serverErr = driver.ActionElement(&ActionNormalParam{
		Element: elementId,
		Text:    "",
	}, Clear)
	if serverErr != nil {
		return
	}
	return
}

// FindInputMethods Get the input-method
func (driver DeviceDriverModel) FindInputMethods() (imeKeyboards []ImeKeyboardModel, err *AppiumError) {
	args := []string{
		"-s",
		driver.DeviceName,
		"shell",
		"ime",
		"list",
		"-s",
	}
	out, err := GetOutPutString(GetAdbPath(), args)
	if err != nil {
		return
	}
	devices := strings.Split(strings.TrimSpace(out), "\n")
	imeKeyboards = sf.Map(devices, func(item string) ImeKeyboardModel {
		return StrConvertImeModel(item)
	})
	fmt.Println(imeKeyboards)
	return
}

// SetKeyboardType Set the keyboard
func (driver DeviceDriverModel) SetKeyboardType(imeKeyboard *ImeKeyboardModel) (err *AppiumError) {
	args := []string{
		"-s",
		driver.DeviceName,
		"shell",
		"ime",
		"enable",
		imeKeyboard.ToString(),
	}
	err = NoOutPutString(GetAdbPath(), args)
	if err != nil {
		return
	}
	args = []string{
		"-s",
		driver.DeviceName,
		"shell",
		"ime",
		"set",
		imeKeyboard.ToString(),
	}
	err = NoOutPutString(GetAdbPath(), args)
	return
}

// GetNetworkStatus
// Use the `adb shell ping` to get the network status
func (driver DeviceDriverModel) GetNetworkStatus() (ret bool, err *AppiumError) {
	args := []string{
		"-s",
		driver.DeviceName,
		"shell",
		"ping",
		"www.baidu.com",
	}
	ret, err = KillLoopCmd(GetAdbPath(), args)
	if err != nil {
		return
	}
	return
}

func (driver DeviceDriverModel) GetElementText(element *FindElementPoint) (value AttributeInterface, elementId string, serverErr *AppiumError) {
	serverErr = driver.ImplicitWait(2 * time.Second)
	if serverErr != nil {
		return
	}
	// 1. Find the element
	elementId, serverErr = driver.FindElement(element)
	if serverErr != nil {
		return
	}

	var result AttributeResponse

	resp, err := driver.Client.R().
		SetSuccessResult(&result).
		Get(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/element/%s/text", driver.Port, driver.SessionId, elementId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Not Found the text",
			ErrorCode: NotFoundAttribute,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Not Found the text",
			ErrorCode: NotFoundAttribute,
		}
		return
	}
	value = &AttributeRetModel{Value: result.Value}
	return
}

func (driver DeviceDriverModel) GetElementTextByElementId(elementId string) (value string, serverErr *AppiumError) {

	var result TextResponseValue

	resp, err := driver.Client.R().
		SetSuccessResult(&result).
		Get(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/element/%s/text", driver.Port, driver.SessionId, elementId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Not Found the text",
			ErrorCode: NotFoundAttribute,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Not Found the text",
			ErrorCode: NotFoundAttribute,
		}
		return
	}
	value = result.Value
	return
}

// GetAllContext
//
//	@Description: Get the context of devices
//	@receiver driver
//	@return ret
//	@return serverErr
func (driver DeviceDriverModel) GetAllContext() (ret []string, serverErr *AppiumError) {
	var response GetContextResponse

	resp, err := driver.Client.R().
		SetSuccessResult(&response).
		Get(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/contexts", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Get context error",
			ErrorCode: GetContextsError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Get ",
			ErrorCode: GetContextsError,
		}
		return
	}
	ret = response.Value
	return
}

func (driver DeviceDriverModel) SetContext(context string) (serverErr *AppiumError) {
	var response SessionResponse

	resp, err := driver.Client.R().
		SetBody(&SetContextParam{
			Name: context,
		}).
		SetSuccessResult(&response).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/context", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Set context error",
			ErrorCode: SetContextError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Set context error",
			ErrorCode: SetContextError,
		}
		return
	}
	return
}

func (driver DeviceDriverModel) PressCode(codeNum int) (serverErr *AppiumError) {
	var response SessionResponse

	resp, err := driver.Client.R().
		SetBody(&PressCodeParam{
			KeyCode: codeNum,
		}).
		SetSuccessResult(&response).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/appium/device/press_keycode", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Android Press Code Error",
			ErrorCode: PressCodeError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Android Press Code Error",
			ErrorCode: PressCodeError,
		}
		return
	}
	return
}

func (driver DeviceDriverModel) OpenAirplaneMode(isOpen bool) (serverErr *AppiumError) {
	//var open int
	var isOpenStr string
	if isOpen {
		//open = 1
		isOpenStr = "disable"
	} else {
		//open = 0
		isOpenStr = "enable"
	}
	args := []string{
		"-s",
		driver.DeviceName,
		"shell",
		"svc",
		"data",
		isOpenStr,
	}
	args2 := []string{
		"-s",
		driver.DeviceName,
		"shell",
		"svc",
		"wifi",
		isOpenStr,
	}
	serverErr = NoOutPutString(GetAdbPath(), args)
	serverErr = NoOutPutString(GetAdbPath(), args2)
	if serverErr != nil {
		return
	}
	args = []string{
		"-s",
		driver.DeviceName,
		"shell",
		"su",
		"-c",
		fmt.Sprintf("'am broadcast -a android.intent.action.AIRPLANE_MODE --ez state %s'", isOpenStr),
	}
	serverErr = NoOutPutString(GetAdbPath(), args)
	return
}

// RemoveApp Remove an app from the device
func (driver DeviceDriverModel) RemoveApp(appName string) (serverErr *AppiumError) {
	var response SessionResponse

	resp, err := driver.Client.R().
		SetBody(&RemoveAppParam{
			BundleId: appName,
		}).
		SetSuccessResult(&response).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/appium/device/remove_app", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Uninstall app failure",
			ErrorCode: RemoveAppError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Uninstall app failure",
			ErrorCode: RemoveAppError,
		}
		return
	}
	return
}

// InstallApp Install the given app onto the device
func (driver DeviceDriverModel) InstallApp(appPath string) (serverErr *AppiumError) {
	var response SessionResponse

	resp, err := driver.Client.R().
		SetBody(&InstallAppParam{
			AppPath: appPath,
		}).
		SetSuccessResult(&response).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/appium/device/install_app", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Install app failure",
			ErrorCode: InstallAppError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Install app failure",
			ErrorCode: InstallAppError,
		}
		return
	}
	return
}

// IsAppInstall Check whether the specified app is installed on the device
func (driver DeviceDriverModel) IsAppInstall(appName string) (ret bool, serverErr *AppiumError) {
	var response TerminateResponse

	resp, err := driver.Client.R().
		SetBody(&RemoveAppParam{
			BundleId: appName,
		}).
		SetSuccessResult(&response).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/appium/device/app_installed", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Not found app",
			ErrorCode: IsAppInstallError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Not found app",
			ErrorCode: IsAppInstallError,
		}
		return
	}
	ret = response.Value
	return
}

// HideKeyboard Hide soft keyboard
func (driver DeviceDriverModel) HideKeyboard() (serverErr *AppiumError) {
	var response SessionResponse

	resp, err := driver.Client.R().
		SetBody(&HideKeyboardParam{
			Strategy: "default",
		}).
		SetSuccessResult(&response).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/appium/device/hide_keyboard", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Not found app",
			ErrorCode: IsAppInstallError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Not found app",
			ErrorCode: IsAppInstallError,
		}
		return
	}
	return
}

// IsKeyboardShown Whether the soft keyboard is shown
func (driver DeviceDriverModel) IsKeyboardShown() (ret bool, serverErr *AppiumError) {
	var response TerminateResponse

	resp, err := driver.Client.R().
		SetSuccessResult(&response).
		Get(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/appium/device/is_keyboard_shown", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Is keyboard shown error",
			ErrorCode: IsKeyboardShowError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Is keyboard shown error",
			ErrorCode: IsKeyboardShowError,
		}
		return
	}
	ret = response.Value
	return
}

// Scroll : scroll on the touch screen using finger based motion events
func (driver DeviceDriverModel) Scroll(start Coordinate, end Coordinate) (serverErr *AppiumError) {
	// Create the parameters
	var actions []ActionRequestChain
	actions = append(actions, ActionRequestChain{
		Type:     "pointerMove",
		Duration: start.GetDuration(),
		X:        start.GetPosition().X,
		Y:        start.GetPosition().Y,
	})
	actions = append(actions, ActionRequestChain{
		Type:   "pointerDown",
		Button: 0,
	})
	actions = append(actions, ActionRequestChain{
		Type:     "pointerMove",
		Duration: end.GetDuration(),
		X:        end.GetPosition().X,
		Y:        end.GetPosition().Y,
		Origin:   "viewport",
	})
	actions = append(actions, ActionRequestChain{
		Type:   "pointerUp",
		Button: 0,
	})
	pointerParam := ActionRequestParams{
		PointerType: "touch",
	}
	actionParam := ActionsRequest{
		Actions:    actions,
		Parameters: pointerParam,
		Id:         "finger1",
		Type:       "pointer",
	}
	var actionTemp []ActionsRequest
	actionTemp = append(actionTemp, actionParam)
	requestParams := &ActionRequestArr{
		Actions: actionTemp,
	}
	var result SessionResponse

	resp, err := driver.Client.R().
		SetBody(requestParams).
		SetSuccessResult(&result).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/actions", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Scroll Action Error",
			ErrorCode: ScrollError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Scroll Action Error",
			ErrorCode: ScrollError,
		}
		return
	}
	return
}

// ActivateApp
//
//	@Description: Activate app by appid
//	@receiver driver
//	@param appId
//	@return serverErr
func (driver DeviceDriverModel) ActivateApp(appId string) (serverErr *AppiumError) {
	var response SessionResponse

	resp, err := driver.Client.R().
		SetBody(&ActivateAppParam{
			AppId: appId,
		}).
		SetSuccessResult(&response).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/appium/device/activate_app", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Activate app failure",
			ErrorCode: ActivateAppError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Activate app failure",
			ErrorCode: ActivateAppError,
		}
		return
	}
	return
}

// GetAvailableIme
//
//	@Description: **Android only** show the available input-methods on the devices
//	@receiver driver
//	@return imeKeyboards
//	@return serverErr
func (driver DeviceDriverModel) GetAvailableIme() (imeKeyboards []string, serverErr *AppiumError) {
	var response GetContextResponse

	resp, err := driver.Client.R().
		SetSuccessResult(&response).
		Get(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/ime/available_engines", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Available ime error",
			ErrorCode: AvailableImeError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Available ime error",
			ErrorCode: AvailableImeError,
		}
		return
	}

	imeKeyboards = response.Value
	return
}

// GetActiveIme
//
//	@Description: Get the active input-method activity and package name
//	@receiver driver
//	@return imeKeyboard
//	@return serverErr
func (driver DeviceDriverModel) GetActiveIme() (imeKeyboard string, serverErr *AppiumError) {
	var response AttributeResponse

	resp, err := driver.Client.R().
		SetSuccessResult(&response).
		Get(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/ime/active_engine", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Available ime error",
			ErrorCode: AvailableImeError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Available ime error",
			ErrorCode: AvailableImeError,
		}
		return
	}

	imeKeyboard = response.Value
	return
}

// ActivateImeBoard
//
//	@Description: Activate the input-method
//	@receiver driver
//	@param ime
//	@return serverErr
func (driver DeviceDriverModel) ActivateImeBoard(ime string) (serverErr *AppiumError) {
	var response SessionResponse

	resp, err := driver.Client.R().
		SetBody(&ImeSearchParam{
			Engine: ime,
		}).
		SetSuccessResult(&response).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/ime/activate", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Available ime error",
			ErrorCode: ActivateImeError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Available ime error",
			ErrorCode: ActivateImeError,
		}
		return
	}
	return
}

// GetCurrentPackage
//
//	@Description: Get the current package name
//	@receiver driver
//	@return packageName
//	@return serverErr
func (driver DeviceDriverModel) GetCurrentPackage() (packageName string, serverErr *AppiumError) {
	var response AttributeResponse

	resp, err := driver.Client.R().
		SetSuccessResult(&response).
		Get(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/appium/device/current_package", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Get package name error",
			ErrorCode: GetPackageNameError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Get package name error",
			ErrorCode: GetPackageNameError,
		}
		return
	}

	packageName = response.Value
	return
}

// ActionChainsMove
//
//	@Description: Slide multi pointers
//	@receiver driver
//	@param coordinates
//	@return serverErr
func (driver DeviceDriverModel) ActionChainsMove(coordinates []Coordinate) (serverErr *AppiumError) {
	var actions []ActionRequestChain
	if len(coordinates) <= 1 {
		serverErr = &AppiumError{
			Message:   "Slide mode just one type",
			ErrorCode: 0,
		}
		return
	}
	actions = append(actions, ActionRequestChain{
		Type:     "pointerMove",
		Duration: 1,
		X:        coordinates[0].GetPosition().X,
		Y:        coordinates[0].GetPosition().Y,
	})
	actions = append(actions, ActionRequestChain{
		Type:   "pointerDown",
		Button: 0,
	})
	tmpCoordinates := coordinates[1:]
	for _, position := range tmpCoordinates {
		actions = append(actions, ActionRequestChain{
			Type:     "pointerMove",
			Duration: position.GetDuration(),
			X:        position.GetPosition().X,
			Y:        position.GetPosition().Y,
		})
	}
	actions = append(actions, ActionRequestChain{
		Type:   "pointerUp",
		Button: 0,
	})
	pointerParam := ActionRequestParams{
		PointerType: "touch",
	}
	actionParam := ActionsRequest{
		Actions:    actions,
		Parameters: pointerParam,
		Id:         "finger1",
		Type:       "pointer",
	}

	var actionTemp []ActionsRequest
	actionTemp = append(actionTemp, actionParam)
	requestParams := &ActionRequestArr{
		Actions: actionTemp,
	}
	var result SessionResponse

	resp, err := driver.Client.R().
		SetBody(requestParams).
		SetSuccessResult(&result).
		Post(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/actions", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Touch Action Error",
			ErrorCode: TouchActionError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Touch Action Error",
			ErrorCode: TouchActionError,
		}
		return
	}
	return
}

// ScreenShot gets a screenshot of the current device state. The returned image is in Base64 encoding.
// Returns the image and an error if failed.
func (driver DeviceDriverModel) ScreenShot() (imageBase64 string, serverErr *AppiumError) {
	var response TextResponseValue

	resp, err := driver.Client.R().
		SetSuccessResult(&response).
		Get(fmt.Sprintf("http://127.0.0.1:%d/wd/hub/session/%s/screenshot", driver.Port, driver.SessionId))
	if err != nil {
		serverErr = &AppiumError{
			Message:   "Get screen shot error",
			ErrorCode: GetScreenShotError,
		}
		return
	}

	if !resp.IsSuccessState() {
		serverErr = &AppiumError{
			Message:   "Get screen shot error",
			ErrorCode: GetScreenShotError,
		}
		return
	}

	imageBase64 = response.Value
	return
}

func (driver DeviceDriverModel) GrantAllPermission(packageName string) (serverErr *AppiumError) {
	args := []string{
		"-s",
		driver.DeviceName,
		"shell",
		"pm",
		"grant",
		packageName,
		"android.permission.WRITE_EXTERNAL_STORAGE",
	}
	serverErr = NoOutPutString(GetAdbPath(), args)
	args = []string{
		"-s",
		driver.DeviceName,
		"shell",
		"pm",
		"grant",
		packageName,
		"android.permission.READ_EXTERNAL_STORAGE",
	}
	serverErr = NoOutPutString(GetAdbPath(), args)
	if serverErr != nil {
		return
	}
	return
}
