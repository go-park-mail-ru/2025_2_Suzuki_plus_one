package yookassa

import (
	"errors"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
)

// See
// - SDK examples
// 		- https://github.com/rvinnie/yookassa-sdk-go/blob/main/docs/examples/01-configuration.md
// - Dashboard
// 		- https://yookassa.ru/my/merchant/integration/http-notifications

type Yookassa struct {
	logger      logger.Logger
	Client      *yookassa.Client
	Handler     *yookassa.PaymentHandler
	redirectURL string
}

func NewYookassa(logger logger.Logger, shopID, secretKey, redirectURL string) (*Yookassa, error) {
	client := yookassa.NewClient(shopID, secretKey)

	// Создаем обработчик настроек
	settingsHandler := yookassa.NewSettingsHandler(client)
	// Получаем информацию о настройках магазина или шлюза
	settings, err := settingsHandler.GetAccountSettings(nil)
	if err != nil {
		return nil, errors.New("NewYookassa: failed to get account settings: " + err.Error())
	}
	logger.Info("Yookassa account settings retrieved",
		logger.ToString("settings.AccountId", settings.AccountId),
		logger.ToAny("settings.Status", settings.Status),
		logger.ToAny("settings.Test", settings.Test),
		logger.ToAny("settings.FiscalizationEnabled", settings.FiscalizationEnabled),
		logger.ToString("settings.ITN", settings.ITN),
		logger.ToInt("settings.PaymentMethodsCount", len(settings.PaymentMethods)),
		logger.ToInt("settings.PayoutMethodsCount", func() int {
			if settings.PayoutMethods != nil {
				return len(*settings.PayoutMethods)
			}
			return 0
		}()),
		logger.ToString("settings.Name", settings.Name),
	)

	return &Yookassa{
		logger:      logger,
		Client:      client,
		Handler:     yookassa.NewPaymentHandler(client),
		redirectURL: redirectURL,
	}, nil
}
