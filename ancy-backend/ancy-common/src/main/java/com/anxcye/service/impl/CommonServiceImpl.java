package com.anxcye.service.impl;


import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.exception.SystemException;
import com.anxcye.service.CommonService;
import com.anxcye.utils.AliOssUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.UUID;

/**
 * @author axy
 * @description
 * @createDate 2024-09-05 11:29:50
 */
@Service
public class CommonServiceImpl implements CommonService {

    @Autowired
    private AliOssUtil aliOssUtil;

    @Override
    public String uploadImg(MultipartFile file) {
        String originalFilename = file.getOriginalFilename();
        String fileName = UUID.randomUUID().toString() + originalFilename.substring(originalFilename.lastIndexOf("."));
        try {
            byte[] fileBytes = file.getBytes();
            return aliOssUtil.upload(fileBytes, fileName);

        } catch (IOException e) {
            throw new SystemException(AppHttpCodeEnum.UPLOAD_ERROR);
        }
    }
}
