package com.anxcye.controller;

import org.springframework.web.bind.annotation.RestController;

@RestController
public class UserController {


//    @Autowired
//    private UserService userService;
//
//    @PostMapping("/login")
//    public ResponseResult<?> login(@RequestBody LoginDto loginDto) {
//        if (StringUtils.isBlank(loginDto.getUserName()) || StringUtils.isBlank(loginDto.getPassword())) {
//            throw new SystemException(AppHttpCodeEnum.REQUIRE_USERNAME);
//        }
//        return ResponseResult.success(userService.login(loginDto));
//    }
//
//    @PostMapping("/logout")
//    public ResponseResult<?> logout() {
//        userService.logout();
//        return ResponseResult.success();
//    }
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
