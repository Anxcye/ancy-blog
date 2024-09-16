package com.anxcye.exception;

import com.anxcye.domain.enums.AppHttpCodeEnum;

public class SystemException extends RuntimeException {
    private AppHttpCodeEnum appHttpCodeEnum;

    public int getCode() {
        return appHttpCodeEnum.getCode();
    }

    public String getMsg() {
        return appHttpCodeEnum.getMsg();
    }

    public AppHttpCodeEnum getAppHttpCodeEnum() {
        return appHttpCodeEnum;
    }

    public SystemException(AppHttpCodeEnum httpCodeEnum) {
        super(httpCodeEnum.getMsg());
        this.appHttpCodeEnum = httpCodeEnum;
    }

}
