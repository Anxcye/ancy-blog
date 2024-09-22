package com.anxcye.service.impl;

import com.alibaba.excel.util.StringUtils;
import com.anxcye.domain.dto.TagDto;
import com.anxcye.domain.dto.TagListDto;
import com.anxcye.domain.entity.Tag;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.TagVo;
import com.anxcye.mapper.TagMapper;
import com.anxcye.service.TagService;
import com.anxcye.utils.BeanCopyUtils;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * @author axy
 * @description 针对表【ancy_tag(标签)】的数据库操作Service实现
 * @createDate 2024-09-20 10:28:04
 */
@Service
public class TagServiceImpl extends ServiceImpl<TagMapper, Tag>
        implements TagService {

    @Override
    public PageResult pageList(TagListDto tagListDto) {
        LambdaQueryWrapper<Tag> tagLambdaQueryWrapper = new LambdaQueryWrapper<>();
        tagLambdaQueryWrapper.like(StringUtils.isNotBlank(tagListDto.getName()),
                Tag::getName, tagListDto.getName());
        tagLambdaQueryWrapper.like(StringUtils.isNotBlank(tagListDto.getRemark()),
                Tag::getRemark, tagListDto.getRemark());

        Page<Tag> tagPage = new Page<>(tagListDto.getPageNum(), tagListDto.getPageSize());
        page(tagPage, tagLambdaQueryWrapper);
        List<TagVo> tagVos = BeanCopyUtils.copyList(tagPage.getRecords(), TagVo.class);
        return new PageResult(tagPage.getTotal(), tagVos);
    }

    @Override
    public boolean addTag(TagDto tagDto) {
        Tag tag = BeanCopyUtils.copyBean(tagDto, Tag.class);
        save(tag);
        return true;
    }

    @Override
    public boolean deleteTag(Long id) {
        removeById(id);
        return true;
    }

    @Override
    public boolean updateTag(Long id, TagDto tagDto) {
        Tag tag = BeanCopyUtils.copyBean(tagDto, Tag.class);
        tag.setId(id);
        updateById(tag);
        return true;
    }

    @Override
    public List<TagVo> selectTagsByArticleId(Long articleId) {

        List<Tag> tags = getBaseMapper().selectTagsByArticleId(articleId);

        return BeanCopyUtils.copyList(tags, TagVo.class);
    }
}




