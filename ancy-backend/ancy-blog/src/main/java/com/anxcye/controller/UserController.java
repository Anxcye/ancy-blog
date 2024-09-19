package com.anxcye.controller;

import com.alibaba.excel.util.StringUtils;
import com.anxcye.domain.dto.LoginDto;
import com.anxcye.domain.dto.RegisterDto;
import com.anxcye.domain.dto.UserDto;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.exception.SystemException;
import com.anxcye.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
public class UserController {


    @Autowired
    private UserService userService;

    @PostMapping("/login")
    public ResponseResult<?> login(@RequestBody LoginDto loginDto) {
        if (StringUtils.isBlank(loginDto.getUserName()) || StringUtils.isBlank(loginDto.getPassword())) {
            throw new SystemException(AppHttpCodeEnum.REQUIRE_USERNAME);
        }
        return ResponseResult.success(userService.login(loginDto));
    }

    @PostMapping("/logout")
    public ResponseResult<?> logout() {
        userService.logout();
        return ResponseResult.success();
    }

    @GetMapping("/info")
    public ResponseResult<?> getUserInfo() {
        return ResponseResult.success(userService.getUserInfo());
    }
    
    @PutMapping("/info")
    public ResponseResult<?> updateUserInfo(@RequestBody UserDto userDto) {
        return ResponseResult.success(userService.updateUserInfo(userDto));
    }

    @PostMapping("/register")
    public ResponseResult<?> register(@RequestBody RegisterDto registerDto) {
        return ResponseResult.success(userService.register(registerDto));
    }
    
}
