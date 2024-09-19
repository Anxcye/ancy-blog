package com.anxcye.domain.dto;

import lombok.Data;

@Data
public class RegisterDto {
    /**
     * 用户名
     */
    private String userName;

    /**
     * 昵称
     */
    private String nickName;

    /**
     * 邮箱
     */
    private String email;

    /**
     * 密码
     */
    private String password;
}
