package com.anxcye.service;

import com.anxcye.domain.dto.TagDto;
import com.anxcye.domain.dto.TagListDto;
import com.anxcye.domain.entity.Tag;
import com.anxcye.domain.result.PageResult;
import com.baomidou.mybatisplus.extension.service.IService;

/**
* @author axy
* @description 针对表【ancy_tag(标签)】的数据库操作Service
* @createDate 2024-09-20 10:28:04
*/
public interface TagService extends IService<Tag> {

    PageResult pageList(TagListDto tagListDto);

    boolean addTag(TagDto tagDto);

    boolean deleteTag(Long id);

    boolean updateTag(Long id, TagDto tagDto);
}
