import moment from 'moment'

export default {
    methods: {
        //时间格式化
        dateFormat(row, column) {
            if (row[column.property] === 0) {
                return '-'
            }
            return moment.unix(row[column.property]).format('YYYY-MM-DD HH:mm:ss');
        }
    }
}
