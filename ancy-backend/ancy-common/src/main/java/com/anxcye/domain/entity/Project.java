package com.anxcye.domain.entity;

import com.baomidou.mybatisplus.annotation.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;
import java.util.Date;

/**
 * project表
 * @TableName ancy_project
 */
@TableName(value ="ancy_project")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Project implements Serializable {
    /**
     * 
     */
    @TableId(type = IdType.AUTO)
    private Long id;

    /**
     * 标题
     */
    private String title;

    /**
     * content
     */
    private String content;

    /**
     * 文章摘要
     */
    private String summary;

    /**
     * 缩略图
     */
    private String thumbnail;

    /**
     * 是否置顶（0否，1是）
     */
    private String isTop;

    /**
     * 状态（0已发布，1草稿）
     */
    private String status;

    /**
     * 0 active 1 archived
     */
    private String type;

    /**
     * 
     */
    private String srcUrl;

    /**
     * 
     */
    private String displayUrl;

    /**
     * 排序
     */
    private Integer orderNum;

    /**
     * 
     */
    private Date beginDate;

    /**
     * 
     */
    @TableField(fill = FieldFill.INSERT)
    private Long createBy;

    /**
     * 
     */
    @TableField(fill = FieldFill.INSERT)
    private Date createTime;

    /**
     * 
     */
    @TableField(fill = FieldFill.UPDATE)
    private Long updateBy;

    /**
     * 
     */
    @TableField(fill = FieldFill.UPDATE)
    private Date updateTime;

    /**
     * 删除标志（0代表未删除，1代表已删除）
     */
    private Integer deleted;

    @TableField(exist = false)
    private static final long serialVersionUID = 1L;
}