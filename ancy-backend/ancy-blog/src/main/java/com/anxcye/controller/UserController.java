package com.anxcye.controller;

import com.alibaba.excel.util.StringUtils;
import com.anxcye.domain.entity.User;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.exception.SystemException;
import com.anxcye.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class UserController {


    @Autowired
    private UserService userService;

    @PostMapping("/login")
    public ResponseResult<?> login(@RequestBody User user) {
        if (StringUtils.isBlank(user.getUserName()) || StringUtils.isBlank(user.getPassword())) {
            throw new SystemException(AppHttpCodeEnum.REQUIRE_USERNAME);
        }
        return ResponseResult.success(userService.login(user));
    }

    @PostMapping("/logout")
    public ResponseResult<?> logout() {
        userService.logout();
        return ResponseResult.success();
    }
}
