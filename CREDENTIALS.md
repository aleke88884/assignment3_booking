# Настройка внешних сервисов для SmartBooking

Этот файл содержит инструкции по настройке API ключей и подключению внешних сервисов для production версии приложения.

## 📁 Содержание

- [Google Cloud Storage](#google-cloud-storage)
- [AWS S3](#aws-s3)
- [Firebase Cloud Messaging (уведомления)](#firebase-cloud-messaging)
- [Email сервис (SendGrid)](#sendgrid-email)
- [SMS сервис (Twilio)](#twilio-sms)
- [Карты (Google Maps / Yandex Maps)](#карты)

---

## 🗄️ Google Cloud Storage

### Шаг 1: Создание проекта

1. Перейдите на [Google Cloud Console](https://console.cloud.google.com/)
2. Создайте новый проект или выберите существующий
3. Включите Google Cloud Storage API

### Шаг 2: Создание Service Account

1. Перейдите в **IAM & Admin** → **Service Accounts**
2. Нажмите **Create Service Account**
3. Заполните:
   - **Name**: `smartbooking-storage`
   - **Role**: `Storage Admin`
4. Создайте JSON ключ и скачайте его

### Шаг 3: Создание bucket

```bash
# Установите gsutil
# Создайте bucket
gsutil mb -l EUROPE-WEST1 gs://smartbooking-photos

# Настройте публичный доступ
gsutil iam ch allUsers:objectViewer gs://smartbooking-photos
```

### Шаг 4: Настройка в .env

```env
STORAGE_TYPE=gcs
STORAGE_BUCKET=smartbooking-photos
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
STORAGE_PUBLIC_URL=https://storage.googleapis.com/smartbooking-photos
```

### Цена

- **Первые 5GB**: Бесплатно
- **Хранение**: ~$0.02/GB в месяц
- **Операции**: ~$0.05 за 10,000 запросов

---

## ☁️ AWS S3

### Шаг 1: Создание IAM пользователя

1. Перейдите в [AWS Console](https://console.aws.amazon.com/)
2. **IAM** → **Users** → **Add User**
3. Выберите **Programmatic access**
4. Прикрепите политику `AmazonS3FullAccess`
5. Сохраните **Access Key ID** и **Secret Access Key**

### Шаг 2: Создание S3 Bucket

1. Перейдите в **S3** → **Create Bucket**
2. Имя: `smartbooking-photos-prod`
3. Регион: `eu-central-1` (Frankfurt)
4. Отключите **Block all public access**
5. Создайте bucket

### Шаг 3: Настройка CORS

В настройках bucket добавьте CORS:

```json
[
    {
        "AllowedHeaders": ["*"],
        "AllowedMethods": ["GET", "PUT", "POST", "DELETE"],
        "AllowedOrigins": ["*"],
        "ExposeHeaders": []
    }
]
```

### Шаг 4: Bucket Policy (публичный доступ к файлам)

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::smartbooking-photos-prod/*"
        }
    ]
}
```

### Шаг 5: Настройка в .env

```env
STORAGE_TYPE=s3
STORAGE_ENDPOINT=
STORAGE_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
STORAGE_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
STORAGE_BUCKET=smartbooking-photos-prod
STORAGE_REGION=eu-central-1
STORAGE_USE_SSL=true
STORAGE_PUBLIC_URL=https://smartbooking-photos-prod.s3.eu-central-1.amazonaws.com
```

### Цена

- **Первые 5GB**: Бесплатно (12 месяцев)
- **Хранение**: ~$0.023/GB в месяц
- **Передача данных**: Первый 1GB бесплатно

---

## 🔔 Firebase Cloud Messaging

### Шаг 1: Создание проекта Firebase

1. Перейдите на [Firebase Console](https://console.firebase.google.com/)
2. **Add Project** → укажите имя `SmartBooking`
3. Выберите **Analytics** (опционально)

### Шаг 2: Получение Server Key

1. Перейдите в **Project Settings** (шестеренка)
2. Вкладка **Cloud Messaging**
3. Скопируйте **Server key**

### Шаг 3: Скачайте Service Account JSON

1. **Project Settings** → **Service Accounts**
2. **Generate new private key**
3. Скачайте JSON файл

### Шаг 4: Web Push Certificates

1. В **Cloud Messaging** → **Web Configuration**
2. Сгенерируйте **Web Push certificates**
3. Скопируйте **Key pair**

### Шаг 5: Настройка в .env

```env
FCM_SERVER_KEY=AAAA1234567890:APA91bGxxx...
FCM_SERVICE_ACCOUNT=/path/to/firebase-service-account.json
FCM_VAPID_KEY=BNxxx...
```

### Пример отправки уведомления (Go)

```go
import "google.golang.org/api/fcm/v1"

func SendNotification(token string, title string, body string) error {
    message := &fcm.Message{
        Token: token,
        Notification: &fcm.Notification{
            Title: title,
            Body:  body,
        },
    }
    // отправка через FCM API
}
```

### Цена

- **Полностью бесплатно** для любого количества уведомлений

---

## 📧 SendGrid Email

### Шаг 1: Регистрация

1. Зарегистрируйтесь на [SendGrid](https://sendgrid.com/)
2. Подтвердите email

### Шаг 2: Создание API Key

1. **Settings** → **API Keys**
2. **Create API Key**
3. Выберите **Full Access**
4. Скопируйте ключ (он показывается только раз!)

### Шаг 3: Настройка отправителя

1. **Settings** → **Sender Authentication**
2. **Verify a Single Sender**
3. Заполните email отправителя (например, `noreply@smartbooking.kz`)

### Шаг 4: Настройка в .env

```env
SENDGRID_API_KEY=SG.xxx...
EMAIL_FROM=noreply@smartbooking.kz
EMAIL_FROM_NAME=SmartBooking
```

### Пример отправки (Go)

```go
import "github.com/sendgrid/sendgrid-go"

func SendEmail(to string, subject string, body string) {
    message := mail.NewSingleEmail(
        mail.NewEmail("SmartBooking", "noreply@smartbooking.kz"),
        subject,
        mail.NewEmail("", to),
        body,
        body,
    )
    client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
    response, err := client.Send(message)
}
```

### Цена

- **100 emails/день**: Бесплатно навсегда
- **40,000 emails/месяц**: $14.95

---

## 📱 Twilio SMS

### Шаг 1: Регистрация

1. Зарегистрируйтесь на [Twilio](https://www.twilio.com/)
2. Подтвердите номер телефона

### Шаг 2: Получение credentials

1. **Console** → **Account Info**
2. Скопируйте:
   - **Account SID**
   - **Auth Token**

### Шаг 3: Получение номера телефона

1. **Phone Numbers** → **Buy a number**
2. Выберите страну (Kazakhstan +7)
3. Купите номер (~$1/месяц)

### Шаг 4: Настройка в .env

```env
TWILIO_ACCOUNT_SID=ACxxx...
TWILIO_AUTH_TOKEN=xxx...
TWILIO_PHONE_NUMBER=+77001234567
```

### Пример отправки SMS (Go)

```go
import "github.com/twilio/twilio-go"

func SendSMS(to string, message string) {
    client := twilio.NewRestClient()
    params := &api.CreateMessageParams{}
    params.SetTo(to)
    params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
    params.SetBody(message)
    
    client.Api.CreateMessage(params)
}
```

### Цена

- **Trial**: $15.50 бесплатных кредитов
- **SMS в Казахстан**: ~$0.08 за сообщение
- **Номер телефона**: $1/месяц

---

## 🗺️ Карты

### Google Maps API

1. [Google Cloud Console](https://console.cloud.google.com/)
2. **APIs & Services** → **Enable APIs**
3. Включите:
   - Maps JavaScript API
   - Geocoding API
   - Places API
4. **Credentials** → **Create Credentials** → **API Key**
5. Ограничьте ключ:
   - **Application restrictions**: HTTP referrers
   - **API restrictions**: выберите нужные API

```env
GOOGLE_MAPS_API_KEY=AIzaSyxxx...
```

**Цена**: 

- $200 бесплатных кредитов в месяц
- ~$7 за 1000 загрузок карты

### Yandex Maps API

1. [Яндекс.Кабинет](https://developer.tech.yandex.ru/)
2. **Получить ключ API**
3. Выберите **JavaScript API**

```env
YANDEX_MAPS_API_KEY=xxx-xxx-xxx
```

**Цена**:

- 25,000 запросов/день бесплатно
- Далее ~$1 за 1000 запросов

---

## 🔐 Безопасность

### Хранение credentials

**НЕ КОММИТЬТЕ** файлы с ключами в Git!

Добавьте в `.gitignore`:

```gitignore
.env
.env.production
*.json
*-key.json
firebase-*.json
credentials/
```

### Использование в production

Используйте переменные окружения или секреты:

**Docker Secrets:**

```yaml
services:
  app:
    environment:
      SENDGRID_API_KEY: ${SENDGRID_API_KEY}
      GOOGLE_MAPS_API_KEY: ${GOOGLE_MAPS_API_KEY}
```

**Kubernetes Secrets:**

```bash
kubectl create secret generic app-secrets \
  --from-literal=sendgrid-key=SG.xxx \
  --from-literal=fcm-key=AAAA...
```

---

## 📋 Checklist для Production

- [ ] AWS S3 / Google Cloud Storage настроен
- [ ] Firebase FCM ключи получены
- [ ] SendGrid API key получен и протестирован
- [ ] Twilio account настроен (опционально)
- [ ] Google Maps / Yandex Maps API ключи получены
- [ ] Все ключи добавлены в `.env` (локально)
- [ ] Все ключи добавлены в secrets manager (production)
- [ ] `.env` файлы добавлены в `.gitignore`
- [ ] CORS настроен для S3
- [ ] Email отправитель верифицирован в SendGrid
- [ ] Протестирована отправка уведомлений

---

## 🆘 Помощь

### Полезные ссылки

- [AWS S3 Documentation](https://docs.aws.amazon.com/s3/)
- [Google Cloud Storage](https://cloud.google.com/storage/docs)
- [Firebase FCM](https://firebase.google.com/docs/cloud-messaging)
- [SendGrid API](https://docs.sendgrid.com/)
- [Twilio Docs](https://www.twilio.com/docs)
- [Google Maps Platform](https://developers.google.com/maps)

### Тестирование локально

Для локальной разработки используйте:

- **Storage**: MinIO (уже настроен в docker-compose)
- **Email**: MailHog или Mailtrap
- **SMS**: Twilio test credentials

---

**Последнее обновление**: 2026-02-05
