package com.anxcye.controller;

import com.anxcye.domain.dto.SettingDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.SettingVo;
import com.anxcye.service.SettingService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/setting")
public class SettingController {

    @Autowired
    private SettingService settingService;

    @PreAuthorize("@ps.hasPermission('content:setting:list')")
    @GetMapping
    public ResponseResult<SettingVo> settingList() {
        return ResponseResult.success(settingService.getBaseSetting());
    }

    @PreAuthorize("@ps.hasPermission('content:setting:edit')")
    @PutMapping
    public ResponseResult<Boolean> settingUpdate(@RequestBody SettingDto settingDto) {
        return ResponseResult.success(settingService.updateBaseSetting(settingDto));
    }
}
