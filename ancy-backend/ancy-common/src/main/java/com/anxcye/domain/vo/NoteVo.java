package com.anxcye.domain.vo;

import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;

import java.io.Serializable;
import java.util.Date;

/**
 * note表
 * @TableName ancy_note
 */
@TableName(value ="ancy_note")
@Data
public class NoteVo implements Serializable {
    /**
     * 
     */
    private Long id;

    /**
     * 文章内容
     */
    private String content;

    /**
     * 是否置顶（0否，1是）
     */
    private String isTop;

    /**
     * 状态（0已发布，1草稿）
     */
    private String status;

    /**
     * 排序
     */
    private Integer orderNum;

    /**
     * 访问量
     */
    private Long viewCount;

    /**
     * 是否允许评论 1是，0否
     */
    private String isComment;

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

    /**
     * 删除标志（0代表未删除，1代表已删除）
     */
    private Integer deleted;
}