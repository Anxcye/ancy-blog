package com.anxcye.controller;

import com.anxcye.service.CommonService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class CommonController {

    @Autowired
    private CommonService commonService;

//    @PostMapping("/upload")
//    public ResponseResult<?> upload(MultipartFile file) {
//        return ResponseResult.success(commonService.uploadImg(file));
//    }
}
