package email

/*
使用示例：

1. 初始化（在 main.go 或启动代码中）：

	import (
		"context"
		"supervillain/internal/common/email"
	)

	// 注册 Email 发送器到 IOC
	emailConfig := &email.Config{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: "your-email@gmail.com",
		Password: "your-app-password", // Gmail 需要使用应用专用密码
		From:     "your-email@gmail.com",
		FromName: "Your App Name",
		TLS:      false, // 587 端口使用 STARTTLS，设置为 false
	}
	email.RegisterSender(emailConfig)

2. 同步发送邮件：

	msg := &email.Message{
		To:      []string{"recipient@example.com"},
		Cc:      []string{"cc@example.com"}, // 可选
		Bcc:     []string{"bcc@example.com"}, // 可选
		Subject: "测试邮件",
		Body:    "这是一封测试邮件",
		HTML:    "<h1>这是一封测试邮件</h1><p>HTML 内容</p>", // 可选，如果提供则优先使用
	}

	err := email.Send(context.Background(), msg)
	if err != nil {
		logger.Error("Failed to send email", "error", err)
	}

3. 异步发送邮件（推荐，使用 WorkerModule）：

	msg := &email.Message{
		To:      []string{"recipient@example.com"},
		Subject: "异步测试邮件",
		Body:    "这是一封异步发送的测试邮件",
	}

	taskID, err := email.SendAsync(context.Background(), msg)
	if err != nil {
		logger.Error("Failed to submit email task", "error", err)
	} else {
		logger.Info("Email task submitted", "task_id", taskID)
	}

	// 可以查询任务状态
	workerModule := ioc.MustGet[*module.WorkerModule]("workerModule")
	taskInfo, exists := workerModule.GetTaskInfo(taskID)
	if exists {
		logger.Info("Email task status", "status", taskInfo.StatusStr)
	}

4. 直接使用发送器（不使用 IOC）：

	sender := email.NewSMTPSender(emailConfig)
	
	// 同步发送
	err := sender.Send(context.Background(), msg)
	
	// 异步发送（需要 WorkerModule 在 IOC 中）
	taskID, err := sender.SendAsync(context.Background(), msg)

5. 常见 SMTP 配置：

	Gmail:
		Host: "smtp.gmail.com"
		Port: 587 (STARTTLS) 或 465 (TLS)
		TLS: false (587) 或 true (465)
		需要启用"应用专用密码"

	QQ 邮箱:
		Host: "smtp.qq.com"
		Port: 587
		TLS: false
		需要开启 SMTP 服务并获取授权码

	163 邮箱:
		Host: "smtp.163.com"
		Port: 25 或 465
		TLS: false (25) 或 true (465)
		需要开启 SMTP 服务并获取授权码

	企业邮箱（以腾讯企业邮箱为例）:
		Host: "smtp.exmail.qq.com"
		Port: 587
		TLS: false
		使用邮箱账号和密码
*/

