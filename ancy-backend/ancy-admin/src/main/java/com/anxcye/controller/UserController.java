package com.anxcye.controller;

import com.anxcye.domain.dto.AdminUserDto;
import com.anxcye.domain.dto.LoginDto;
import com.anxcye.domain.dto.UserListDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.UserService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/user")
public class UserController {

    @Autowired
    private UserService userService;

    @PostMapping("/login")
    public ResponseResult<?> login(@RequestBody LoginDto loginDto) {
        return ResponseResult.success(userService.adminLogin(loginDto));
    }

    @GetMapping("/routers")
    public ResponseResult<?> getRouters() {
        return ResponseResult.success(userService.getRouters());
    }

    @PostMapping("/logout")
    public ResponseResult<?> logout() {
        userService.adminLogout();
        return ResponseResult.success();
    }

    @GetMapping("/page")
    public ResponseResult<?> getUserPage(@ParameterObject UserListDto userListDto) {
        return ResponseResult.success(userService.getPage(userListDto));
    }

    @GetMapping("/{id}")
    public ResponseResult<?> getUserById(@PathVariable Long id) {
        return ResponseResult.success(userService.getAdminById(id));
    }

    @PostMapping
    public ResponseResult<?> addUser(@RequestBody AdminUserDto adminUserDto) {
        return ResponseResult.success(userService.addAdmin(adminUserDto));
    }

    @PutMapping("/{id}")
    public ResponseResult<?> updateUser(@PathVariable Long id, @RequestBody AdminUserDto adminUserDto) {
        return ResponseResult.success(userService.updateAdmin(id, adminUserDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<?> deleteUser(@PathVariable Long id) {
        return ResponseResult.success(userService.deleteUser(id));
    }
}
