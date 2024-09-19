package com.anxcye.domain.enums;

import lombok.Getter;

@Getter
public enum AppHttpCodeEnum {
    // 成功
    SUCCESS(200, "操作成功"),
    // 登录
    NEED_LOGIN(401, "需要登录后操作"),
    NO_OPERATOR_AUTH(403, "无权限操作"),
    SYSTEM_ERROR(500, "出现错误"),
    USERNAME_EXIST(501, "用户名已存在"),
    PHONE_NUMBER_EXIST(502, "手机号已存在"),
    EMAIL_EXIST(503, "邮箱已存在"),
    REQUIRE_USERNAME(504, "请提供用户名和密码"),
    LOGIN_ERROR(505, "用户名密码错误或用户被禁用"),
    AUTH_ERROR(506, "认证失败"),
    CONTENT_NOT_NULL(507, "内容不能为空"),
    UPLOAD_ERROR(508, "上传失败"),
    USERINFO_NOT_NULL(509, "用户信息不完整"),
    TOKEN_INVALID(510, "无效的token");

    final int code;
    final String msg;

    AppHttpCodeEnum(int code, String errorMessage) {
        this.code = code;
        this.msg = errorMessage;
    }
}


