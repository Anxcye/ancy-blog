package com.anxcye.domain.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class AddArticleDto {
    // 标题
    private String title;
    // 文章内容
    private String content;
    // 文章摘要
    private String summary;
    // 所属分类id
    private Long categoryId;
    // 缩略图
    private String thumbnail;
    // 是否置顶（0否，1是）
    private String isTop;
    // 状态（0已发布，1草稿）
    private String status;
    // 是否允许评论 1是，0否
    private String isComment;
    // 标签ID列表
    private List<Long> tags;

}
