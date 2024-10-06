package com.anxcye.domain.vo;

import lombok.Data;

import java.util.Date;
import java.util.List;

@Data
public class ArticleCardVo {
    /**
     *
     */
    private Long id;

    /**
     * 标题
     */
    private String title;

    /**
     * 文章摘要
     */
    private String summary;

    /**
     * 所属分类名
     */
    private String categoryName;

    /**
     * 所属分类id
     */
    private Long categoryId;

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
     * 文章类型 0普通文章 1首页文章 2关联项目 3友链文章
     */
    private Integer type;

    /**
     * 排序
     */
    private Integer orderNum;

    /**
     * tags
     */
    private List<TagVo> tags;
    /**
     * 访问量
     */
    private Long viewCount;

    /**
     *
     */
    private Long createBy;

    /**
     *
     */
    private Date createTime;

    private Date updateTime;

}
