package com.anxcye.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.Tag;
import com.anxcye.service.TagService;
import com.anxcye.mapper.TagMapper;
import org.springframework.stereotype.Service;

/**
* @author axy
* @description 针对表【ancy_tag(标签)】的数据库操作Service实现
* @createDate 2024-09-20 10:28:04
*/
@Service
public class TagServiceImpl extends ServiceImpl<TagMapper, Tag>
    implements TagService{

}




