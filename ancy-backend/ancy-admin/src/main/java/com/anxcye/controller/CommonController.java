package com.anxcye.controller;

import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.CommonService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

@RestController
public class CommonController {

    @Autowired
    private CommonService commonService;

    @PreAuthorize("@ps.hasPermission('content:article:upload')")
    @PostMapping("/upload")
    public ResponseResult<?> upload(MultipartFile file) {
        return ResponseResult.success(commonService.uploadImg(file));
    }
}
