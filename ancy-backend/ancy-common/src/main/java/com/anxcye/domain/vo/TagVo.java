package com.anxcye.domain.vo;

import lombok.Data;

import java.util.Date;

@Data
public class TagVo {
       /**
     *
     */
    private Long id;

    /**
     * 标签名
     */
    private String name;

    /**
     * 备注
     */
    private String remark;

    /**
     *
     */
    private Long createBy;

    /**
     *
     */
    private Date createTime;

    /**
     *
     */
    private Long updateBy;

    /**
     *
     */
    private Date updateTime;

}
