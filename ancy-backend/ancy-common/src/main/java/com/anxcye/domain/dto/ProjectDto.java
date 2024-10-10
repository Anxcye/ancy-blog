package com.anxcye.domain.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Date;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class ProjectDto {
    /**
     * 标题
     */
    private String title;

    /**
     * content
     */
    private String content;

    /**
     * 摘要
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
}
