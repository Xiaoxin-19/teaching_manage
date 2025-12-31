export interface SaveRecordData {
  student_id: number;
  student_name?: string;
  teaching_date: string;
  start_time: string;
  end_time: string;
  remark: string;
}

export interface StudentOption {
  id: number;
  name: string;
}
