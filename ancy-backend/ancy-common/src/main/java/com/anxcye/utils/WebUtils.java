package com.anxcye.utils;


import com.alibaba.fastjson.JSON;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.result.ResponseResult;
import jakarta.servlet.http.HttpServletResponse;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;

public class WebUtils {
    /**
     * 将字符串渲染到客户端
     *
     * @param response 渲染对象
     * @param appHttpCodeEnum   待渲染的字符串
     */
    public static void renderString(HttpServletResponse response, AppHttpCodeEnum appHttpCodeEnum) {
        try {
            String result = JSON.toJSONString(ResponseResult.error(appHttpCodeEnum));
            response.setStatus(appHttpCodeEnum.getCode());
            response.setContentType("application/json");
            response.setCharacterEncoding("utf-8");
            response.getWriter().print(result);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public static void setDownLoadHeader(String filename, HttpServletResponse response) throws UnsupportedEncodingException {
        response.setContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet");
        response.setCharacterEncoding("utf-8");
        String fname = URLEncoder.encode(filename, StandardCharsets.UTF_8).replaceAll("\\+", "%20");
        response.setHeader("Content-disposition", "attachment; filename=" + fname);
    }
}