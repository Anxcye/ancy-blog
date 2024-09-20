package com.anxcye.controller;

import com.anxcye.domain.dto.LoginDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/user")
public class UserController {

    @Autowired
    private UserService userService;

    @PostMapping("/login")
    public ResponseResult<?>login(@RequestBody LoginDto loginDto){
        return  ResponseResult.success(userService.adminLogin(loginDto));
    }

    @GetMapping("/routers")
    public ResponseResult<?>  getRouters(){
        return ResponseResult.success(userService.getRouters());
    }
}
