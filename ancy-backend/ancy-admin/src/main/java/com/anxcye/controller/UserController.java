package com.anxcye.controller;

import com.anxcye.domain.dto.AdminUserDto;
import com.anxcye.domain.dto.LoginDto;
import com.anxcye.domain.dto.UserListDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.AdminUserVo;
import com.anxcye.domain.vo.RouterVo;
import com.anxcye.domain.vo.UserVo;
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
    public ResponseResult<AdminUserVo> login(@RequestBody LoginDto loginDto) {
        return ResponseResult.success(userService.adminLogin(loginDto));
    }

    @GetMapping("/routers")
    public ResponseResult<RouterVo> getRouters() {
        return ResponseResult.success(userService.getRouters());
    }

    @PostMapping("/logout")
    public ResponseResult<Boolean> logout() {

        return ResponseResult.success(userService.adminLogout());
    }

    @GetMapping("/page")
    public ResponseResult<PageResult> userPage(@ParameterObject UserListDto userListDto) {
        return ResponseResult.success(userService.getPage(userListDto));
    }

    @GetMapping("/{id}")
    public ResponseResult<UserVo> userGetById(@PathVariable Long id) {
        return ResponseResult.success(userService.getAdminById(id));
    }

    @PostMapping
    public ResponseResult<Long> userAdd(@RequestBody AdminUserDto adminUserDto) {
        return ResponseResult.success(userService.addAdmin(adminUserDto));
    }

    @PutMapping("/{id}")
    public ResponseResult<Boolean> userUpdate(@PathVariable Long id, @RequestBody AdminUserDto adminUserDto) {
        return ResponseResult.success(userService.updateAdmin(id, adminUserDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> userDelete(@PathVariable Long id) {
        return ResponseResult.success(userService.deleteUser(id));
    }
}
