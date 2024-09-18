package com.anxcye.service;

import org.springframework.web.multipart.MultipartFile;

/**
* @author axy
* @description
* @createDate 2024-09-12 13:57:34
*/
public interface CommonService  {

    String uploadImg(MultipartFile file);
}
