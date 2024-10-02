package com.anxcye.controller;

import com.anxcye.domain.dto.RoleDto;
import com.anxcye.domain.dto.RoleListDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.RoleVo;
import com.anxcye.service.RoleService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/role")
public class RoleController {
    @Autowired
    private RoleService roleService;

    @GetMapping("/page")
    public ResponseResult<PageResult> rolePage(@ParameterObject RoleListDto roleListDto) {
        return ResponseResult.success(roleService.getPage(roleListDto));
    }

    @GetMapping("/list")
    public ResponseResult<List<RoleVo>> roleList() {
        return ResponseResult.success(roleService.getList());
    }

    @GetMapping("/{id}")
    public ResponseResult<RoleVo> roleGetById(@PathVariable Long id) {
        return ResponseResult.success(roleService.getRoleById(id));
    }

    @PostMapping
    public ResponseResult<Boolean> roleAdd(@RequestBody RoleDto roleDto) {
        return ResponseResult.success(roleService.addRole(roleDto));
    }

    @PutMapping("/{id}")
    public ResponseResult<Boolean> roleUpdate(@PathVariable Long id, @RequestBody RoleDto roleDto) {
        return ResponseResult.success(roleService.updateRole(id, roleDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> roleDelete(@PathVariable Long id) {
        return ResponseResult.success(roleService.deleteRole(id));
    }
}
