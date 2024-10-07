package com.anxcye.domain.vo;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;
import java.util.Date;


@Data
@AllArgsConstructor
@NoArgsConstructor
public class TimelineVo implements Serializable {
    /**
     * ID
     */
    private Long id;

    /**
     * 操作时间
     */
    private Date operateTime;

    /**
     *
     */
    private String summary;

    /**
     * 操作的方法名
     */
    private String methodName;

    /**
     * 返回值
     */
    private String returnValue;

}