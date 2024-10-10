package com.anxcye.domain.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class ReadDto {
    /**
     * 出处
     */
    private String source;

    /**
     * 内容
     */
    private String content;

    /**
     * 作者
     */
    private String author;

    /**
     * 添加来源 0 手动添加 1安读
     */
    private Integer addFrom;
}
