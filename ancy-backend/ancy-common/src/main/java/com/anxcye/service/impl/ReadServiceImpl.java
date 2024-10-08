package com.anxcye.service.impl;

import com.anxcye.domain.dto.ReadPageDto;
import com.anxcye.domain.entity.Read;
import com.anxcye.domain.result.PageResult;
import com.anxcye.mapper.ReadMapper;
import com.anxcye.service.ReadService;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

/**
* @author axy
* @description 针对表【ancy_read(阅读表)】的数据库操作Service实现
* @createDate 2024-10-06 22:25:31
*/
@Service
public class ReadServiceImpl extends ServiceImpl<ReadMapper, Read>
    implements ReadService{

    @Override
    public PageResult getReadPage(ReadPageDto readPageDto) {
        LambdaQueryWrapper<Read> wrapper = new LambdaQueryWrapper<>();
        wrapper.like(StringUtils.hasText(readPageDto.getSource()), Read::getSource, readPageDto.getSource());
        wrapper.like(StringUtils.hasText(readPageDto.getContent()), Read::getContent, readPageDto.getContent());
        wrapper.like(StringUtils.hasText(readPageDto.getAuthor()), Read::getAuthor, readPageDto.getAuthor());
        wrapper.orderByDesc(Read::getCreateTime);
        Page<Read> page = new Page<>(readPageDto.getPageNum(), readPageDto.getPageSize());
        page(page, wrapper);
        return new PageResult(page.getTotal(), page.getRecords());
    }
}



