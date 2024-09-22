package com.anxcye.controller;

import com.anxcye.domain.dto.MenuDto;
import com.anxcye.domain.dto.MenuListDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.MenuService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/menus")
public class MenuController {

    @Autowired
    private MenuService menuService;

    @GetMapping("/list")
    public ResponseResult<?> list(@ParameterObject MenuListDto menuListDto) {
        return ResponseResult.success(menuService.listMenus(menuListDto));
    }

    @PostMapping
    public ResponseResult<?> addMenu(@RequestBody MenuDto menuDto) {
        return ResponseResult.success(menuService.addMenu(menuDto));
    }

    @PutMapping("/{id}")
    public ResponseResult<?> updateMenu(@PathVariable Long id, @RequestBody MenuDto menuDto) {
        return ResponseResult.success(menuService.updateMenu(id, menuDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<?> deleteMenu(@PathVariable Long id) {
        return ResponseResult.success(menuService.deleteMenu(id));
    }
    

}
