package com.anxcye.controller;

import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.SettingService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;


@RestController
@RequestMapping("/setting")
public class SettingController {

    @Autowired
    private SettingService settingService;

    @GetMapping("/base")
    public ResponseResult<?> settingGetBase() {
        return ResponseResult.success(settingService.getBaseSetting());
    }

//    @PutMapping
//    public ResponseResult<?> settingUpdateBase(@RequestBody SettingDto settingDto) {
//        return ResponseResult.success(settingService.updateBaseSetting(settingDto));
//    }
//



}
