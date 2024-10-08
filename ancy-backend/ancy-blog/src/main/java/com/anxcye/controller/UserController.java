package com.anxcye.controller;

import com.anxcye.domain.dto.LoginDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.AdminUserVo;
import com.anxcye.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/user")
public class UserController {


    @Autowired
    private UserService userService;

    @PostMapping("/login")
    public ResponseResult<AdminUserVo> login(@RequestBody LoginDto loginDto) {
        return ResponseResult.success(userService.adminLogin(loginDto));
    }

    @PostMapping("/logout")
    public ResponseResult<Boolean> logout() {
        return ResponseResult.success(userService.adminLogout());
    }
//
//    @GetMapping("/info")
//    public ResponseResult<?> getUserInfo() {
//        return ResponseResult.success(userService.getUserInfo());
//    }
//
//    @PutMapping("/info")
//    public ResponseResult<?> updateUserInfo(@RequestBody UserDto userDto) {
//        return ResponseResult.success(userService.updateUserInfo(userDto));
//    }
//
//    @PostMapping("/register")
//    public ResponseResult<?> register(@RequestBody RegisterDto registerDto) {
//        return ResponseResult.success(userService.register(registerDto));
//    }
    
}
