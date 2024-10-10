package com.anxcye.service;

import com.anxcye.domain.dto.ReadDto;
import com.anxcye.domain.dto.ReadPageDto;
import com.anxcye.domain.entity.Read;
import com.anxcye.domain.result.PageResult;
import com.baomidou.mybatisplus.extension.service.IService;

/**
* @author axy
* @description 针对表【ancy_read(阅读表)】的数据库操作Service
* @createDate 2024-10-06 22:25:31
*/
public interface ReadService extends IService<Read> {

    PageResult getReadPage(ReadPageDto readPageDto);

    Long addRead(ReadDto readDto);

    Boolean updateRead(Long id, ReadDto readDto);

    Boolean deleteRead(Long id);
}
