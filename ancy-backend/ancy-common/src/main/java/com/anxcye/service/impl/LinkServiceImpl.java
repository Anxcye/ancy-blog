package com.anxcye.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.Link;
import com.anxcye.service.LinkService;
import com.anxcye.mapper.LinkMapper;
import org.springframework.stereotype.Service;

/**
* @author axy
* @description 针对表【ancy_link(友链)】的数据库操作Service实现
* @createDate 2024-09-10 20:33:41
*/
@Service
public class LinkServiceImpl extends ServiceImpl<LinkMapper, Link>
    implements LinkService{

}




