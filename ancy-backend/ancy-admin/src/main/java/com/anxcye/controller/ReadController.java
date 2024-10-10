package com.anxcye.controller;

import com.anxcye.domain.dto.ReadDto;
import com.anxcye.domain.dto.ReadPageDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.ReadService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;
@RestController
@RequestMapping("/read")
public class ReadController {

    @Autowired
    private ReadService readService;

    @PreAuthorize("@ps.hasPermission('content:read:list')")
    @GetMapping
    public ResponseResult<PageResult> readPage(@ParameterObject ReadPageDto readPageDto) {
        return ResponseResult.success(readService.getReadPage(readPageDto));
    }

    @PreAuthorize("@ps.hasPermission('content:read:add')")
    @PostMapping
    public ResponseResult<Long> readAdd(@RequestBody ReadDto readDto) {
        return ResponseResult.success(readService.addRead(readDto));
    }

    @PreAuthorize("@ps.hasPermission('content:read:edit')")
    @PutMapping("/{id}")
    public ResponseResult<Boolean> readUpdate(@PathVariable Long id, @RequestBody ReadDto readDto) {
        return ResponseResult.success(readService.updateRead(id, readDto));
    }

    @PreAuthorize("@ps.hasPermission('content:read:remove')")
    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> readDelete(@PathVariable Long id) {
        return ResponseResult.success(readService.deleteRead(id));
    }
}
