package appium_cli

import "github.com/imroc/req/v3"

type PlatformType int64

// Appium platforms support Android, IOS, Mac, Windows...
const (
	Android PlatformType = 0
	IOS                  = 1
	Mac                  = 2
	Windows              = 3
)

type ActionType int64

// Different Action have different actions
const (
	Click ActionType = iota
	SendKeys
	Press
	Clear
)

type DeviceInfoResponse struct {
	DeviceScreenSize string `json:"deviceScreenSize"`
}

type SessionResponse struct {
	SessionId string             `json:"sessionId"`
	Status    int64              `json:"status"`
	Value     DeviceInfoResponse `json:"value,omitempty"`
}

type AttributeResponse struct {
	SessionId string `json:"sessionId"`
	Status    int64  `json:"status"`
	Value     string `json:"value,omitempty"`
}

type TextResponseValue struct {
	SessionId string `json:"sessionId"`
	Value     string `json:"value,omitempty"`
}

type ValueErrorMsg struct {
	Message string `json:"message"`
}

type SessionErrorResponse struct {
	SessionResponse
	Value ValueErrorMsg `json:"value"`
}

// DesiredCapabilities 实际上是 W3C 协议中的 "alwaysMatch" 部分
type DesiredCapabilities struct {
	PlatformName          string `json:"platformName"` // 标准字段，无前缀
	PlatformVersion       string `json:"appium:platformVersion"`
	DeviceName            string `json:"appium:deviceName"`
	AppPackage            string `json:"appium:appPackage"`
	AppActivity           string `json:"appium:appActivity"`
	NewCommandTimeout     uint64 `json:"appium:newCommandTimeout"`
	AndroidInstallTimeout uint64 `json:"appium:androidInstallTimeout"`
	AutomationName        string `json:"appium:automationName"`
	SystemPort            uint64 `json:"appium:systemPort"`
	Udid                  string `json:"appium:udid"`
	NoReset               bool   `json:"appium:noReset"`
	// 如果有其他字段，也需要加 appium: 前缀
}

// DeviceCapabilityModel user start the device capability
type DeviceCapabilityModel struct {
	Platform              PlatformType
	PlatformVersion       string `json:"platformVersion"`
	DeviceName            string `json:"deviceName"`
	AppPackage            string `json:"appPackage"`
	AppActivity           string `json:"appActivity"`
	NewCommandTimeout     uint64 `json:"newCommandTimeout"`
	AndroidInstallTimeout uint64 `json:"androidInstallTimeout"`
	AutomationName        string `json:"automationName"`
	SystemPort            uint64 `json:"systemPort"`
	Udid                  string `json:"udid"`
	NoReset               bool   `json:"noReset"`
	ChromedriverExcutable string `json:"chromedriverExcutable"`
	Port                  int64
	Client                *req.Client
}

type CapabilitiesWrapper struct {
	AlwaysMatch DesiredCapabilities `json:"alwaysMatch"`
}

type AppiumParameter struct {
	Capabilities CapabilitiesWrapper `json:"capabilities"`
}

type DeviceDriverModel struct {
	SessionId  string
	Client     *req.Client
	Port       int64
	DeviceName string
}

// AppiumBy
// @Description: Access
// Accessibility ID	Read a unique identifier for a UI element. For XCUITest it is the element's accessibility-id attribute. For Android it is the element's content-desc attribute.
// Class name	For IOS it is the full name of the XCUI element and begins with XCUIElementType. For Android it is the full name of the UIAutomator2 class (e.g.: android.widget.TextView)
// ID	Native element identifier. resource-id for android; name for iOS.
type AppiumBy int64

const (
	AccessibilityID AppiumBy = iota
	ID
	XPath
	UiSelector
)

// FindElementParam request post the parameters
type FindElementParam struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

type FindElementPoint struct {
	AppiumBy
	Value string
}

// GetUsingType
//
//	@Description: Get the using type
//	@receiver model
//	@return string
func (model *FindElementPoint) GetUsingType() string {
	if model.AppiumBy == AccessibilityID {
		return "accessibility id"
	} else if model.AppiumBy == ID {
		return "id"
	} else if model.AppiumBy == XPath {
		return "xpath"
	} else {
		return "-android uiautomator"
	}
}

// ElementResponse get the element response
type ElementResponse struct {
	SessionId string            `json:"sessionId"`
	Status    int64             `json:"status"`
	Value     ElementValueModel `json:"value"`
}

type ElementsResponse struct {
	SessionId string              `json:"sessionId"`
	Status    int64               `json:"status"`
	Value     []ElementValueModel `json:"value"`
}

type ElementValueModel struct {
	ELEMENT string `json:"ELEMENT"`
}

type ActionRequestParam struct {
	Element string `json:"element"`
}

type SendKeysParam struct {
	Text string `json:"text"`
}

type ActionNormalParam struct {
	Element string
	Text    string
}

type ActionChainParams struct {
	X int64
	Y int64
}

type ActionRequestChain struct {
	Type     string `json:"type"`
	Duration int64  `json:"duration,omitempty"`
	X        int64  `json:"x,omitempty"`
	Y        int64  `json:"y,omitempty"`
	Button   int64  `json:"button,omitempty"`
	Origin   string `json:"origin,omitempty"`
}

type ActionRequestParams struct {
	PointerType string `json:"pointerType"`
}

type ActionsRequest struct {
	Actions    []ActionRequestChain `json:"actions"`
	Parameters ActionRequestParams  `json:"parameters"`
	Id         string               `json:"id"`
	Type       string               `json:"type"`
}

type ActionRequestArr struct {
	Actions []ActionsRequest `json:"actions"`
}

// StartActivityParam start the another app activity params
type StartActivityParam struct {
	AppPackage  string `json:"appPackage"`
	AppActivity string `json:"appActivity"`
}

type ImplicitWaitParam struct {
	Seconds int `json:"ms"`
}

// ElementAttributeType the attribute of element
type ElementAttributeType int64

const (
	Checked ElementAttributeType = iota
	Clickable
	Enabled
	Displayed
	Text
	Selected
)

type AttributeModel struct {
	AttType ElementAttributeType
}

type AppPropParam struct {
	AppId string `json:"appId"`
}

// GetAttributeStr
//
//	@Description: Get the using type
//	@receiver model
//	@return string
func (model *AttributeModel) GetAttributeStr() string {
	if model.AttType == Checked {
		return "checked"
	} else if model.AttType == Clickable {
		return "clickable"
	} else if model.AttType == Enabled {
		return "enabled"
	} else if model.AttType == Displayed {
		return "displayed"
	} else if model.AttType == Text {
		return "text"
	} else {
		return "selected"
	}
}

type TerminateResponse struct {
	Value     bool   `json:"value"`
	SessionId string `json:"sessionId"`
	Status    int    `json:"status"`
}

type GetContextResponse struct {
	Value     []string `json:"value"`
	SessionId string   `json:"sessionId"`
	Status    int      `json:"status"`
}

type SetContextParam struct {
	Name string `json:"name"`
}

type PressCodeParam struct {
	KeyCode int `json:"keycode"`
}

type RemoveAppParam struct {
	BundleId string `json:"bundleId"`
}

type InstallAppParam struct {
	AppPath string `json:"appPath"`
}

type HideKeyboardParam struct {
	Strategy string `json:"strategy"`
}

type ActivateAppParam struct {
	AppId string `json:"appId"`
}

type ImeSearchParam struct {
	Engine string `json:"engine"`
}
