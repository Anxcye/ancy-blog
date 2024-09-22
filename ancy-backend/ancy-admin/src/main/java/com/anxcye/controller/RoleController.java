package com.anxcye.controller;

import com.anxcye.domain.dto.RoleDto;
import com.anxcye.domain.dto.RoleListDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.RoleService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/role")
public class RoleController {
    @Autowired
    private RoleService roleService;

    @GetMapping("/page")
    public ResponseResult<?> page(@ParameterObject RoleListDto roleListDto) {
        return ResponseResult.success(roleService.getPage(roleListDto));
    }

    @GetMapping("/{id}")
    public ResponseResult<?> getById(@PathVariable Long id) {
        return ResponseResult.success(roleService.getRoleById(id));
    }

    @PostMapping
    public ResponseResult<?> add(@RequestBody RoleDto roleDto) {
        return ResponseResult.success(roleService.addRole(roleDto));
    }

    @PutMapping("/{id}")
    public ResponseResult<?> update(@PathVariable Long id, @RequestBody RoleDto roleDto) {
        return ResponseResult.success(roleService.updateRole(id, roleDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<?> delete(@PathVariable Long id) {
        return ResponseResult.success(roleService.deleteRole(id));
    }
}
