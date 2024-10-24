package com.anxcye.controller;

import com.anxcye.domain.dto.MenuDto;
import com.anxcye.domain.dto.MenuListDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.MenuVo;
import com.anxcye.service.MenuService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/menus")
public class MenuController {

    @Autowired
    private MenuService menuService;

    @PreAuthorize("@ps.hasPermission('system:menu:list')")
    @GetMapping("/list")
    public ResponseResult<List<MenuVo>> menuList(@ParameterObject MenuListDto menuListDto) {
        return ResponseResult.success(menuService.listMenus(menuListDto));
    }

    @PreAuthorize("@ps.hasPermission('system:menu:list')")
    @GetMapping("/tree")
    public ResponseResult<List<MenuVo>> menuTree() {
        return ResponseResult.success(menuService.treeMenus());
    }

    @PreAuthorize("@ps.hasPermission('system:menu:list')")
    @GetMapping("/list/role/{roleId}")
    public ResponseResult<List<MenuVo>> menuListByRoleId(@PathVariable Long roleId) {
        return ResponseResult.success(menuService.selectMenuByRoleId(roleId));
    }

    @PreAuthorize("@ps.hasPermission('system:menu:list')")
    @GetMapping("/page")
    public ResponseResult<PageResult> menuPage(@ParameterObject MenuListDto menuListDto) {
        return ResponseResult.success(menuService.pageMenus(menuListDto));
    }

    @PreAuthorize("@ps.hasPermission('system:menu:add')")
    @PostMapping
    public ResponseResult<Boolean> menuAdd(@RequestBody MenuDto menuDto) {
        return ResponseResult.success(menuService.addMenu(menuDto));
    }

    @PreAuthorize("@ps.hasPermission('system:menu:edit')")
    @PutMapping("/{id}")
    public ResponseResult<Boolean> menuUpdate(@PathVariable Long id, @RequestBody MenuDto menuDto) {
        return ResponseResult.success(menuService.updateMenu(id, menuDto));
    }

    @PreAuthorize("@ps.hasPermission('system:menu:remove')")
    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> menuDelete(@PathVariable Long id) {
        return ResponseResult.success(menuService.deleteMenu(id));
    }
}
