package com.anxcye.controller;

import com.anxcye.domain.entity.User;
import com.anxcye.domain.result.ResponseResult;
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
        return ResponseResult.success(userService.login(user));
    }
}
