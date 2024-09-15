package com.anxcye.domain.result;

import lombok.Data;

import java.io.Serializable;

@Data
public class ResponseResult<T> implements Serializable {

    private Integer code;
    private String msg;
    private T data;

    public static <T> ResponseResult<T> success() {
        ResponseResult<T> result = new ResponseResult<>();
        result.code = 200;
        return result;
    }

    public static <T> ResponseResult<T> success(T object) {
        ResponseResult<T> result = new ResponseResult<>();
        result.data = object;
        result.code = 200;
        return result;
    }

    public static <T> ResponseResult<T> error(String msg) {
        ResponseResult<T> result = new ResponseResult<>();
        result.msg = msg;
        result.code = 0;
        return result;
    }

}
